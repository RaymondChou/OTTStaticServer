#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <string.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/stat.h>
#include "fdfs_client.h"
#include "logger.h"

typedef struct {
    char *msg; //当成功的时候,是返回图片的id,example:"group1/M00/00/00/wKgBP1NxvSqH9qNuAAAED6CzHYE179.jpg",当失败的时候是返回错误消息
    int  result; //１表示成功,０表示失败
}responseData;


responseData upload_file(char *conf_filename, char *local_filename);
responseData delete_file(char *conf_filename, char *file_name);
