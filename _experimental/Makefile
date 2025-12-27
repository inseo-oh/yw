PROJECTDIR := $(CURDIR)

################################################################################
# Platform configuration (such as target name, initial CFLAFGS, etc...)
################################################################################

HOSTCC := $(CC)
ifeq ($(PLATFORM),wasm)
# Emscripten (WASM) ############################################################
CC     := emcc
OUTDIR := $(PROJECTDIR)/out_wasm
TARGET := $(OUTDIR)/yw.html
else ifeq ($(PLATFORM),aos68k)
# AmigaOS 680x0 ################################################################
CC     := vc
OUTDIR := $(PROJECTDIR)/out_aos68k
TARGET := $(OUTDIR)/yw
CFLAGS += +aos68k -lauto -c99
else ifeq ($(PLATFORM),aosppc)
# AmigaOS PowerPC ##############################################################
CC     := vc
OUTDIR := $(PROJECTDIR)/out_aosppc
TARGET := $(OUTDIR)/yw
CFLAGS += +aosppc -lauto -c99
else
# Linux ########################################################################
CFLAGS += -fsanitize=address -g
OUTDIR := $(PROJECTDIR)/out_linux
TARGET := $(OUTDIR)/yw
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

ifneq ($(CC),vc)
CFLAGS += -std=c99
CFLAGS += -Wall -Wextra -Werror -pedantic
CFLAGS += -Wno-error=unused-function -Wno-error=unused-const-variable
endif
CFLAGS += -O3
CFLAGS += -I$(PROJECTDIR)/include

################################################################################
# Build rules
################################################################################

all: $(TARGET)

all-platforms:
	make all
	make all PLATFORM=wasm
	make all PLATFORM=aos68k
	make all PLATFORM=aosppc

clean:
	@rm -f $(OBJS) $(DEPS) $(TARGET)

clean-all-platforms:
	make clean
	make clean PLATFORM=wasm
	make clean PLATFORM=aos68k
	make clean PLATFORM=aosppc

cleandep:
	@rm -f $(DEPS)

prepare:
	@mkdir -p $(OBJDIRS)

$(OBJDIR)/%.o: %.c
	$(info [Target C]    $@)
ifeq ($(CC),vc)
	@$(HOSTCC) -MM -MF $(patsubst %.o, %.d, $@) $<
	@$(CC) -o $@ -c $< $(CFLAGS)
else
	@$(CC) -o $@ -c -MMD $< $(CFLAGS)
endif

$(TARGET): prepare $(OBJS)
	$(info [Target EXE]  $(TARGET))
	@$(CC) -o $(TARGET) $(OBJS) $(CFLAGS)

-include $(DEPS)

.PHONY: all prepare clean cleandep
