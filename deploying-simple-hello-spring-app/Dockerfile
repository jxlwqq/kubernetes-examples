FROM maven:3-openjdk-8 AS compile_stage
WORKDIR /app

COPY .mvn .mvn
COPY pom.xml .
COPY src src

RUN mvn clean package -U -DskipTests

FROM openjdk:8-jdk-alpine
ENV PROJECT_NAME hello-spring
ENV PROJECT_VERSION 0.0.1-SNAPSHOT
WORKDIR /app

COPY --from=compile_stage /app/target/${PROJECT_NAME}-${PROJECT_VERSION}.jar .

CMD ["sh", "-c", "java -jar /app/${PROJECT_NAME}-${PROJECT_VERSION}.jar"]
