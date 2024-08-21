package io.exsql.platform.sql.engine.volcano.operators

import io.exsql.platform.sql.engine.Row

import java.util.function.Predicate
import scala.compiletime.uninitialized

class Filter(input: Operator, predicate: Predicate[Row]) extends Operator {
  private var row: Row = uninitialized

  override def open(): Unit = input.open()
  override def close(): Unit = input.close()
  override def current(): Row = row
  override def next(): Boolean = {
    while input.next() do {
      row = input.current()
      if predicate.test(row) then return true
    }

    false
  }
}
