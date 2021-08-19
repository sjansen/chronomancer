docker-compose up

poetry install
./manage.py migrate
./manage.py createsuperuser
open http://localhost:8000/admin/
