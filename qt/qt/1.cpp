#include <QApplication>
#include <QPushButton>

// #include <QtWidgets/QApplication>
// #include <QtWidgets/QPushButton>
// https://blog.csdn.net/friendbkf/article/details/45440175

int main(int argc, char* argv[]){
	QApplication app(argc,argv);
	QPushButton pushButton(QObject::tr("Hello Qt!"));
	pushButton.show();
	QObject::connect(&pushButton, SIGNAL(clicked()), &app, SLOT(quit()));
	return app.exec();
}
