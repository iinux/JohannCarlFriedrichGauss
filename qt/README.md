qmake -project
qmake qt.pro
make

if qt5 add:
QT += core gui
greaterThan(QT_MAJOR_VERSION, 4): QT += widgets
