package io.exsql.platform.sql

import net.sf.jsqlparser.parser.CCJSqlParserManager
import net.sf.jsqlparser.statement.Statement
import org.slf4j.LoggerFactory

import java.io.StringReader

object Parser {

  private val logger = LoggerFactory.getLogger("io.exsql.platform.sql.SqlParser")

  private val parserManager = CCJSqlParserManager()

  def parse(sql: String): Statement = {
    logger.debug(s"Processing sql query:\n $sql")
    parserManager.parse(StringReader(sql))
  }

}
