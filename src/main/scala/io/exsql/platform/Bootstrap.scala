package io.exsql.platform

import io.exsql.platform.sql.Parser
import io.exsql.platform.sql.engine.Engine.{CloseableIterator, RowBatch}
import io.exsql.platform.sql.engine.{Engine, Row}

object Bootstrap extends App {
  val statement = Parser.parse(
    """
      |SELECT id, name
      |FROM users
      |WHERE salary > 4500
    """.stripMargin
  )

  val partitions = Engine.execute(statement, partitions = new CloseableIterator[RowBatch] {
    private var consumed = 0

    override def close(): Unit = ()
    override def hasNext: Boolean = {
      if consumed == 0 then {
        consumed += 1
        return true
      }

      false
    }

    override def next(): RowBatch = new RowBatch:
      private val rows = Iterator(
        Row(IArray(1, "John", 5000)),
        Row(IArray(2, "Jane", 6000)),
        Row(IArray(3, "Doe", 4000)),
      )

      override def close(): Unit = ()
      override def hasNext: Boolean = rows.hasNext
      override def next(): Row = rows.next()
  })

  partitions.foreach { partition =>
    partition.foreach(row => println(row.toString))
  }
}
