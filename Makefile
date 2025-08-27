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
# Tue, 26 Aug 2025 11:35:01 +0200
# Joe
#
# GNU Makefile

TARGET			:= hf
SHELL			:= /bin/sh
SRC_DIR			:= ./src/
SRC_NAME		:= *.go
MAN_DIR			:= ./man/
MAN_SRC			:= ${TARGET}.1
SRC				:= $(addprefix ${SRC_DIR}, ${SRC_NAME})
DEST			:= /usr
.DEFAULT_GOAL	:= ${TARGET}

run: ${SRC}
	go run ${SRC_DIR}

${TARGET}: ${SRC}
	go build -o ${TARGET} ${SRC_DIR}

install:
	mkdir -p ${DEST}/bin
	cp -f ${TARGET} ${DEST}/bin/hf
	mkdir -p ${DEST}/share/man/man1
	gzip -k ${MAN_DIR}/${MAN_SRC}
	mv -f ${MAN_DIR}/${MAN_SRC}.gz ${DEST}/share/man/man1/${MAN_SRC}.gz

uninstall:
	rm -f ${DEST}/bin/hf
	rm -f ${DEST}/share/man/man1/${MAN_SRC}.gz

release: ${SRC}
	gzip -k ${MAN_DIR}/${MAN_SRC}
	mv -f ${MAN_DIR}/${MAN_SRC}.gz .
	GOOS=darwin GOARCH=arm64   go build -o ${TARGET} ${SRC_DIR}
	tar -zcf hf_v1.0_darwin_arm64.tar.gz  ${TARGET} ${MAN_SRC}.gz README.md LICENSE
	rm -f ${TARGET}
	GOOS=darwin GOARCH=amd64   go build -o ${TARGET} ${SRC_DIR}
	tar -zcf hf_v1.0_darwin_x86_64.tar.gz  ${TARGET} ${MAN_SRC}.gz README.md LICENSE
	rm -f ${TARGET}
	GOOS=freebsd GOARCH=arm64  go build -o ${TARGET} ${SRC_DIR}
	tar -zcf hf_v1.0_freebsd_arm64.tar.gz  ${TARGET} ${MAN_SRC}.gz README.md LICENSE
	rm -f ${TARGET}
	GOOS=freebsd GOARCH=armv6  go build -o ${TARGET} ${SRC_DIR}
	tar -zcf hf_v1.0_freebsd_armv6.tar.gz  ${TARGET} ${MAN_SRC}.gz README.md LICENSE
	rm -f ${TARGET}
	GOOS=freebsd GOARCH=amd64  go build -o ${TARGET} ${SRC_DIR}
	tar -zcf hf_v1.0_freebsd_x86_64.tar.gz ${TARGET} ${MAN_SRC}.gz README.md LICENSE
	rm -f ${TARGET}
	GOOS=linux GOARCH=arm64    go build -o ${TARGET} ${SRC_DIR}
	tar -zcf hf_v1.0_linux_arm64.tar.gz    ${TARGET} ${MAN_SRC}.gz README.md LICENSE
	rm -f ${TARGET}
	GOOS=linux GOARCH=armv6    go build -o ${TARGET} ${SRC_DIR}
	tar -zcf hf_v1.0_linux_armv6.tar.gz    ${TARGET} ${MAN_SRC}.gz README.md LICENSE
	rm -f ${TARGET}
	GOOS=linux GOARCH=amd64    go build -o ${TARGET} ${SRC_DIR}
	tar -zcf hf_v1.0_linux_x86_64.tar.gz   ${TARGET} ${MAN_SRC}.gz README.md LICENSE
	rm -f ${TARGET}
	rm -f ${MAN_SRC}.gz

clean:
	go clean
	rm -f ${TARGET} *.gz

.PHONY: hf run clean install uninstall
