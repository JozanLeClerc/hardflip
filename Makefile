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
SRC				 = $(addprefix ${SRC_DIR}, ${SRC_NAME})
DEST			:= /usr/local/bin
.DEFAULT_GOAL	:= ${TARGET}

run: ${SRC}
	go run ${SRC_DIR}

${TARGET}: ${SRC}
	go build -o ${TARGET} ${SRC_DIR}

install: ${TARGET}
	mkdir -p ${DEST}
	cp -f ${TARGET} ${DEST}
	# man shit
	# mkdir -p $(DESTDIR)$(MANPREFIX)/man1
	# cp -f man/lowbat.1 $(DESTDIR)$(MANPREFIX)/man1/lowbat.1

clean:
	go clean
	rm -f ${TARGET}

.PHONY: hf run clean
