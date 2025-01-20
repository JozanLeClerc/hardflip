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
MAN_DIR			:= ./man/
SRC				:= $(addprefix ${SRC_DIR}, ${SRC_NAME})
DEST			:= /usr
.DEFAULT_GOAL	:= ${TARGET}

run: ${SRC}
	go run -tags debug ${SRC_DIR}

${TARGET}: ${SRC}
	go build -o ${TARGET} ${SRC_DIR}

install:
	mkdir -p ${DEST}/bin
	cp -f ${TARGET} ${DEST}/bin/hf
	mkdir -p ${DEST}/share/man/man1
	gzip ${MAN_DIR}/hf.1
	cp -f man/hf.1.gz ${DEST}/share/man/man1/hf.1.gz
	gzip -d ${MAN_DIR}/hf.1.gz

uninstall:
	rm -f ${DEST}/bin/hf
	rm -f ${DEST}/share/man/man1/hf.1.gz

clean:
	go clean
	rm -f ${TARGET}

.PHONY: hf run clean install uninstall
