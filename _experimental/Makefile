
PROJECTDIR := $(CURDIR)
OUTDIR      = $(PROJECTDIR)/out
OBJDIR      = $(OUTDIR)/obj

SRCS      := $(shell find . -name '*.c')
OBJ_NAMES  = $(patsubst %.c, %.o, $(SRCS))
TARGET     = $(OUTDIR)/yw

OBJS      = $(addprefix $(OBJDIR)/, $(OBJ_NAMES))
DEPS      = $(patsubst %.o, %.d, $(OBJS))
OBJDIRS   = $(sort $(dir $(OBJS)))

CFLAGS  = -std=c99
CFLAGS += -Wall -Wextra -Werror -pedantic
CFLAGS += -I$(PROJECTDIR)/include
CFLAGS += -Wno-error=unused-function -Wno-error=unused-const-variable -g 
CFLAGS += -fsanitize=undefined
CFLAGS += -O3

all: $(TARGET)

clean:
	@rm -f $(OBJS) $(DEPS) $(TARGET)

cleandep:
	@rm -f $(DEPS)

prepare:
	@mkdir -p $(OBJDIRS)

$(OBJDIR)/%.o: %.c
	$(info [Target C]    $@)
	@$(CC) -o $@ -c -MMD $< $(CFLAGS)

$(TARGET): prepare $(OBJS)
	$(info [Target EXE]  $(TARGET))
	@$(CC) -o $(TARGET) $(OBJS) $(CFLAGS)

-include $(DEPS)

.PHONY: all prepare clean cleandep