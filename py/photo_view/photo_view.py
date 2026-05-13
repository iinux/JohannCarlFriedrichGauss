from flask import Flask

import views_index
import views_rtsp
import views_upload
import views_words
from paths import upload_dir, work_dir

app = Flask(__name__, static_folder=work_dir + '/static', static_url_path='')
app.config['UPLOAD_FOLDER'] = upload_dir
app.config['MAX_CONTENT_LENGTH'] = 500 * 1024 * 1024  # 500MB max file size
app.secret_key = 'upload-secret-key-for-flash-messages'

app.register_blueprint(views_index.bp)
app.register_blueprint(views_words.bp)
app.register_blueprint(views_rtsp.bp)
app.register_blueprint(views_upload.bp)


if __name__ == "__main__":
    #app.run(host='192.168.1.4')
    app.run()
