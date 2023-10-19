protobuf
syntax = "proto3";

package lesson;

message CreateLessonRequest {
  string text = 1;
  string translation = 2;
}

message CreateLessonResponse {
  string id = 1;
}

service LessonService {
  rpc CreateLesson(CreateLessonRequest) returns (CreateLessonResponse) {}
}
