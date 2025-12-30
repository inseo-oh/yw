ifeq ($(PROJECTDIR),)
$(error PROJECTDIR must be set)
endif

ifeq ($(TARGET_NAME),)
$(error TARGET_NAME must be set)
endif

################################################################################
# Platform configuration (such as target name, initial CFLAFGS, etc...)
################################################################################

HOSTCC := $(CC)
ifeq ($(PLATFORM),wasm)
# Emscripten (WASM) ############################################################
CC     := emcc
OUTDIR := $(PROJECTDIR)/out_wasm
TARGET := $(OUTDIR)/$(TARGET_NAME).html
ifeq ($(DEBUG),true)
CFLAGS += -fsanitize=address -g
endif
else ifeq ($(PLATFORM),aos68k)
# AmigaOS 680x0 ################################################################
CC     := vc
OUTDIR := $(PROJECTDIR)/out_aos68k
TARGET := $(OUTDIR)/$(TARGET_NAME)
CFLAGS += +aos68k -lauto -c99 -lmieee
else ifeq ($(PLATFORM),aosppc)
# AmigaOS PowerPC ##############################################################
CC     := vc
OUTDIR := $(PROJECTDIR)/out_aosppc
TARGET := $(OUTDIR)/$(TARGET_NAME)
CFLAGS += +aosppc -lauto -c99 -lm
else
# Linux ########################################################################
ifeq ($(DEBUG),true)
CFLAGS += -fsanitize=address -g
endif
OUTDIR := $(PROJECTDIR)/out_linux
TARGET := $(OUTDIR)/$(TARGET_NAME)
endif

################################################################################
# Set up source/object file list and rest of CFLAGS
################################################################################

OBJDIR     = $(OUTDIR)/obj
SRCS      := $(shell find . -name '*.c') $(ADD_SRCS)
OBJ_NAMES  = $(patsubst %.c, %.o, $(abspath $(SRCS)))

OBJS      = $(addprefix $(OBJDIR)/, $(OBJ_NAMES))
DEPS      = $(patsubst %.o, %.d, $(OBJS))
OBJDIRS   = $(sort $(dir $(OBJS)))

ifneq ($(CC),vc)
CFLAGS += -std=c99
CFLAGS += -Wall -Wextra -Werror -pedantic
CFLAGS += -Wno-error=unused-function -Wno-error=unused-const-variable
else
# We urn off some warnings:
# - 51  (bitfield type non-portable)
# - 214 (suspicious format string)
# Both of these seems to have false positives on vbcc when boolean type is used. 
CFLAGS += -dontwarn=51,214
endif
CFLAGS += -I$(PROJECTDIR)/include

ifneq ($(DEBUG),true)
CFLAGS += -O3
endif

################################################################################
# Build rules
################################################################################

all: $(TARGET)

all-platforms:
	$(MAKE) all
	$(MAKE) all PLATFORM=wasm
	$(MAKE) all PLATFORM=aos68k
	$(MAKE) all PLATFORM=aosppc

clean:
	@rm -f $(OBJS) $(DEPS) $(TARGET)

clean-all-platforms:
	$(MAKE) clean
	$(MAKE) clean PLATFORM=wasm
	$(MAKE) clean PLATFORM=aos68k
	$(MAKE) clean PLATFORM=aosppc

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

.PHONY: all all-platforms clean clean-all-platforms prepare cleandep
