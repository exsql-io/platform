package io.exsql.platform.sql.engine.volcano.operators
import io.exsql.platform.sql.engine.Row

import scala.compiletime.uninitialized

class Project(input: Operator, projected: IArray[Int]) extends Operator {
  private var row: Row = uninitialized

  override def open(): Unit = input.open()
  override def close(): Unit = input.close()
  override def current(): Row = project(row)
  override def next(): Boolean = {
    if input.next() then {
      row = input.current()
      return true
    }

    false
  }

  private def project(row: Row): Row = {
    Row(projected.map(row.get))
  }
}
