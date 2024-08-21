package io.exsql.platform.sql.engine

import io.exsql.platform.sql.engine.volcano.operators.{Filter, Project, Scan}
import net.sf.jsqlparser.expression.{Expression, LongValue}
import net.sf.jsqlparser.expression.operators.relational.GreaterThan
import net.sf.jsqlparser.schema.Column
import net.sf.jsqlparser.statement.Statement
import net.sf.jsqlparser.statement.select.{Select, SelectItem}

import java.util.function.Predicate
import scala.collection.mutable
import scala.jdk.CollectionConverters.ListHasAsScala

object Engine {
  trait CloseableIterator[T] extends Iterator[T] {
    def close(): Unit
  }

  trait RowBatch extends CloseableIterator[Row]

  def execute(statement: Statement, partitions: CloseableIterator[RowBatch]): Iterator[RowBatch] = {
    statement match
      case select: Select =>
        val plain = select.getPlainSelect
        val predicate = createPredicate(plain.getWhere())
        val projections = createProjections(plain.getSelectItems.asScala)

        partitions.map { partition =>
          val scan = Scan(partition)
          val filter = Filter(scan, predicate)
          val project = Project(filter, projections)

          new RowBatch:
            override def close(): Unit = project.close()
            override def hasNext: Boolean = project.next()
            override def next(): Row = project.current()
        }
  }

  private def createPredicate(where: Expression): Predicate[Row] = {
    where match
      case gt: GreaterThan =>
        (row: Row) =>
          val left = row.getAs[Int](columnIndexOf(gt.getLeftExpression().asInstanceOf[Column].getColumnName))
          val right = gt.getRightExpression.asInstanceOf[LongValue].getValue
          left > right
  }

  private def createProjections(items: mutable.Buffer[SelectItem[?]]): IArray[Int] = {
    IArray(items.map { item =>
      item.getExpression() match
        case column: Column => columnIndexOf(column.getColumnName)
    }.toSeq: _*)
  }

  private def columnIndexOf(colum: String): Int = {
    colum match
      case "id" => 0
      case "name" => 1
      case "salary" => 2
  }
}
