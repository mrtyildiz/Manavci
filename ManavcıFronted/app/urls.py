# urunler/urls.py
from django.urls import path
from .views import *

urlpatterns = [
    path('', index, name='index'),
    path('index/', index, name='index'),
    path('register/', register, name='register'),
    path('login/', login, name='login'),
    path('logout/', logout_view, name='logout'),
    path('send_mail/', send_mail, name='send_mail'),
]
