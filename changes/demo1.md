# Changes

Code differences compared to source project.

## buf.gen.config.yaml (+1 -1)

```diff
@@ -2,6 +2,6 @@
 inputs:
   - directory: internal
 plugins:
-  - local: ["go", "run", "google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11"]
+  - local: ["go", "run", "google.golang.org/protobuf/cmd/protoc-gen-go@latest"]
     out: internal
     opt: paths=source_relative
```

## buf.gen.yaml (+3 -3)

```diff
@@ -2,10 +2,10 @@
 inputs:
   - directory: api
 plugins:
-  - local: ["go", "run", "google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11"]
+  - local: ["go", "run", "google.golang.org/protobuf/cmd/protoc-gen-go@latest"]
     out: api
     opt: paths=source_relative
-  - local: ["go", "run", "google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.2"]
+  - local: ["go", "run", "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"]
     out: api
     opt: paths=source_relative
   - local: ["go", "run", "github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v3@latest"]
@@ -14,7 +14,7 @@
   - local: ["go", "run", "github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v3@latest"]
     out: api
     opt: paths=source_relative
-  - local: ["go", "run", "github.com/google/gnostic/cmd/protoc-gen-openapi@v0.7.1"]
+  - local: ["go", "run", "github.com/google/gnostic/cmd/protoc-gen-openapi@latest"]
     out: .
     strategy: all
     opt:
```

## buf.yaml (+1 -1)

```diff
@@ -2,6 +2,6 @@
 modules:
   - path: api
   - path: internal
-  - path: third_party
+  - path: proto3ps
 deps:
   - buf.build/googleapis/googleapis
```

## internal/biz/student.go (+5 -9)

```diff
@@ -131,17 +131,13 @@
 }
 
 func (uc *StudentUsecase) ListStudents(ctx context.Context, page int32, pageSize int32) ([]*Student, int32, *ebzkratos.Ebz) {
-	if page < 1 {
-		page = 1
-	}
-	if pageSize < 1 {
-		pageSize = 10
-	}
+	must.True(page >= 1)
+	must.True(pageSize >= 1)
 
 	db := uc.data.DB().WithContext(ctx)
 
-	var total int64
-	if err := db.Model(&Student{}).Count(&total).Error; err != nil {
+	var count int64
+	if err := db.Model(&Student{}).Count(&count).Error; err != nil {
 		return nil, 0, ebzkratos.New(pb.ErrorDbError("count students: %v", err))
 	}
 
@@ -149,5 +145,5 @@
 	if err := db.Order("id").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&items).Error; err != nil {
 		return nil, 0, ebzkratos.New(pb.ErrorDbError("list students: %v", err))
 	}
-	return items, int32(total), nil
+	return items, int32(count), nil
 }
```

## internal/service/student.go (+6 -0)

```diff
@@ -73,6 +73,12 @@
 }
 
 func (s *StudentService) ListStudents(ctx context.Context, req *pb.ListStudentsRequest) (*pb.ListStudentsReply, error) {
+	if req.Page < 1 {
+		return nil, pb.ErrorBadParam("PAGE MUST BE POSITIVE")
+	}
+	if req.PageSize < 1 {
+		return nil, pb.ErrorBadParam("PAGE_SIZE MUST BE POSITIVE")
+	}
 	students, count, ebz := s.uc.ListStudents(ctx, req.Page, req.PageSize)
 	if ebz != nil {
 		return nil, ebz.Erk
```

## proto3ps/errors/errors.proto (+18 -0)

```diff
@@ -0,0 +1,18 @@
+syntax = "proto3";
+
+package errors;
+
+option go_package = "github.com/go-kratos/kratos/v3/errors;errors";
+option java_multiple_files = true;
+option java_package = "com.github.kratos.errors";
+option objc_class_prefix = "KratosErrors";
+
+import "google/protobuf/descriptor.proto";
+
+extend google.protobuf.EnumOptions {
+  int32 default_code = 1108;
+}
+
+extend google.protobuf.EnumValueOptions {
+  int32 code = 1109;
+}
```

## third_party/errors/errors.proto (+0 -18)

```diff
@@ -1,18 +0,0 @@
-syntax = "proto3";
-
-package errors;
-
-option go_package = "github.com/go-kratos/kratos/v3/errors;errors";
-option java_multiple_files = true;
-option java_package = "com.github.kratos.errors";
-option objc_class_prefix = "KratosErrors";
-
-import "google/protobuf/descriptor.proto";
-
-extend google.protobuf.EnumOptions {
-  int32 default_code = 1108;
-}
-
-extend google.protobuf.EnumValueOptions {
-  int32 code = 1109;
-}
```

