package io.exsql.platform.sql.engine

trait Row extends Any {
  def get(ordinal: Int): Any
  def getAs[T](ordinal: Int): T
}

object Row {
  def apply(array: IArray[Any]): Row = ArrayRow(array)

  private[engine] class ArrayRow(array: IArray[Any]) extends AnyVal with Row {
    override def get(ordinal: Int): Any = array(ordinal)

    override def getAs[T](ordinal: Int): T = {
      array(ordinal).asInstanceOf[T]
    }

    override def toString: String = {
      array.mkString("[", ",", "]")
    }
  }
}