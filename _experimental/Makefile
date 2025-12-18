
PROJECTDIR := $(CURDIR)

################################################################################
# Platform configuration (such as target name, initial CFLAFGS, etc...)
################################################################################

ifeq ($(PLATFORM),wasm)
CC      = emcc
OUTDIR  = $(PROJECTDIR)/out_wasm
TARGET  = $(OUTDIR)/yw.html
else
CFLAGS += -fsanitize=undefined -g
OUTDIR  = $(PROJECTDIR)/out_linux
TARGET  = $(OUTDIR)/yw
endif

################################################################################
# Set up source/object file list and rest of CFLAGS
################################################################################

OBJDIR     = $(OUTDIR)/obj
SRCS      := $(shell find . -name '*.c')
OBJ_NAMES  = $(patsubst %.c, %.o, $(SRCS))

OBJS      = $(addprefix $(OBJDIR)/, $(OBJ_NAMES))
DEPS      = $(patsubst %.o, %.d, $(OBJS))
OBJDIRS   = $(sort $(dir $(OBJS)))

CFLAGS += -std=c99
CFLAGS += -Wall -Wextra -Werror -pedantic
CFLAGS += -I$(PROJECTDIR)/include
CFLAGS += -Wno-error=unused-function -Wno-error=unused-const-variable
CFLAGS += -O3

################################################################################
# Build rules
################################################################################

all: $(TARGET)

all-platforms:
	make all
	make all PLATFORM=wasm

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
