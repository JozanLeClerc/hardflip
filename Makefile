# ========================
# =====    ===============
# ======  ================
# ======  ================
# ======  ====   ====   ==
# ======  ===     ==  =  =
# ======  ===  =  ==     =
# =  ===  ===  =  ==  ====
# =  ===  ===  =  ==  =  =
# ==     =====   ====   ==
# ========================
#
# hardflip: Makefile
# Tue Jan 23 11:16:43 2024
# Joe
#
# GNU Makefile

TARGET			:= hf
SHELL			:= /bin/sh
SRC_DIR			:= ./src/
SRC_NAME		:= *.go
CONF_DIR		:= ./src/
SRC				:= $(addprefix ${SRC_DIR}, ${SRC_NAME})
DEST			:= /usr
XDG_CONFIG_HOME ?= $(HOME)/.config
.DEFAULT_GOAL	:= ${TARGET}

run: ${SRC}
	go run ${SRC_DIR}

${TARGET}: ${SRC}
	go build -o ${TARGET} ${SRC_DIR}

install:
	# mkdir -p ${DEST}/bin
	# cp -f ${TARGET} ${DEST}/bin
	# man shit
	# mkdir -p $(DESTDIR)/share/man/man1
	# cp -f man/lowbat.1 $(DESTDIR)/share/man/man1/lowbat.1
	# mkdir -p $(XDG_CONFIG_HOME)/hf
	# cp -f $(CONF_DIR)/* $(XDG_CONFIG_HOME)/hf
	# cp -f $(CONF_DIR)/config.sample.yml $(XDG_CONFIG_HOME)/config.yml

clean:
	go clean
	rm -f ${TARGET}

.PHONY: hf run clean
