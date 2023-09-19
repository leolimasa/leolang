; unnamed_addr is essentially private to this module
@.str = private unnamed_addr constant [11 x i8] c"Holy shit\0", align 1

define i32 @main() local_unnamed_addr #0 {
  %1 = tail call i32 (i32, i8*, ...) @__printf_chkk(i32 1, i8* getelementptr inbounds ([11 x i8], [11 x i8]* @.str, i64 0, i64 0)) #2
  ret i32 0
}


declare i32 @__printf_chkk(i32, i8*, ...) local_unnamed_addr #1


