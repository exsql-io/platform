package io.exsql.platform.sql.engine.volcano.operators

import io.exsql.platform.sql.engine.Engine.RowBatch
import io.exsql.platform.sql.engine.{Row}

class Scan(batch: RowBatch) extends Operator {
  override def open(): Unit = ()
  override def next(): Boolean = batch.hasNext
  override def close(): Unit = batch.close()
  override def current(): Row = batch.next()
}
