package io.exsql.platform.sql.engine.volcano.operators

import io.exsql.platform.sql.engine.Row

trait Operator {
  def open(): Unit
  def next(): Boolean
  def close(): Unit
  def current(): Row
}
