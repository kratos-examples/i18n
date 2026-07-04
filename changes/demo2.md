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

## internal/biz/article.go (+10 -18)

```diff
@@ -145,17 +145,13 @@
 }
 
 func (uc *ArticleUsecase) ListArticles(ctx context.Context, page int32, pageSize int32) ([]*Article, int32, *ebzkratos.Ebz) {
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
-	if err := db.Model(&Article{}).Count(&total).Error; err != nil {
+	var count int64
+	if err := db.Model(&Article{}).Count(&count).Error; err != nil {
 		return nil, 0, ebzkratos.New(pb.ErrorDbError("count articles: %v", err))
 	}
 
@@ -163,7 +159,7 @@
 	if err := db.Order("id").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&items).Error; err != nil {
 		return nil, 0, ebzkratos.New(pb.ErrorDbError("list articles: %v", err))
 	}
-	return items, int32(total), nil
+	return items, int32(count), nil
 }
 
 // ListStudentArticles returns one student's articles, one page at a time. The
@@ -174,17 +170,13 @@
 // 而不是往 ListArticles 上塞过滤参数。
 func (uc *ArticleUsecase) ListStudentArticles(ctx context.Context, studentID int64, page int32, pageSize int32) ([]*Article, int32, *ebzkratos.Ebz) {
 	must.True(studentID > 0)
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
-	if err := db.Model(&Article{}).Where("student_id = ?", studentID).Count(&total).Error; err != nil {
+	var count int64
+	if err := db.Model(&Article{}).Where("student_id = ?", studentID).Count(&count).Error; err != nil {
 		return nil, 0, ebzkratos.New(pb.ErrorDbError("count student articles: %v", err))
 	}
 
@@ -192,5 +184,5 @@
 	if err := db.Where("student_id = ?", studentID).Order("id").Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&items).Error; err != nil {
 		return nil, 0, ebzkratos.New(pb.ErrorDbError("list student articles: %v", err))
 	}
-	return items, int32(total), nil
+	return items, int32(count), nil
 }
```

## internal/service/article.go (+12 -0)

```diff
@@ -79,6 +79,12 @@
 }
 
 func (s *ArticleService) ListArticles(ctx context.Context, req *pb.ListArticlesRequest) (*pb.ListArticlesReply, error) {
+	if req.Page < 1 {
+		return nil, pb.ErrorBadParam("PAGE MUST BE POSITIVE")
+	}
+	if req.PageSize < 1 {
+		return nil, pb.ErrorBadParam("PAGE_SIZE MUST BE POSITIVE")
+	}
 	articles, count, ebz := s.uc.ListArticles(ctx, req.Page, req.PageSize)
 	if ebz != nil {
 		return nil, ebz.Erk
@@ -93,6 +99,12 @@
 func (s *ArticleService) ListStudentArticles(ctx context.Context, req *pb.ListStudentArticlesRequest) (*pb.ListArticlesReply, error) {
 	if req.StudentId <= 0 {
 		return nil, pb.ErrorBadParam("STUDENT_ID IS REQUIRED")
+	}
+	if req.Page < 1 {
+		return nil, pb.ErrorBadParam("PAGE MUST BE POSITIVE")
+	}
+	if req.PageSize < 1 {
+		return nil, pb.ErrorBadParam("PAGE_SIZE MUST BE POSITIVE")
 	}
 	articles, count, ebz := s.uc.ListStudentArticles(ctx, req.StudentId, req.Page, req.PageSize)
 	if ebz != nil {
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

