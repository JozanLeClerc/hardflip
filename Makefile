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
.DEFAULT_GOAL	:= ${TARGET}

run: ${SRC}
	go run ${SRC_DIR} -qwe

${TARGET}: ${SRC}
	go build -o ${TARGET} ${SRC_DIR}

install:
	mkdir -p ${DEST}/bin
	cp -f ${TARGET} ${DEST}/bin
	# man shit
	# mkdir -p $(DESTDIR)/share/man/man1
	# cp -f man/lowbat.1 $(DESTDIR)/share/man/man1/lowbat.1

clean:
	go clean
	rm -f ${TARGET}

.PHONY: hf run clean
