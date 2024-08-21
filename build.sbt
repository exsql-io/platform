ThisBuild / organization := "io.exsql"
ThisBuild / scalaVersion := "3.4.2"
ThisBuild / version      := "0.0.1"

lazy val root = (project in file("."))
  .settings(
    name := "platform",
    libraryDependencies ++= Seq(
      "com.github.jsqlparser" % "jsqlparser" % "5.0",
      "org.slf4j" % "slf4j-api" % "2.0.16",
      "ch.qos.logback" % "logback-classic" % "1.5.7"
    )
  )
