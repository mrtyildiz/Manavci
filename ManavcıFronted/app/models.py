from django.db import models


class ManavciUser(models.Model):
    first_name = models.CharField(max_length=100)
    last_name = models.CharField(max_length=100)
    email = models.EmailField(unique=True)
    gender = models.CharField(max_length=10)
    date_of_birth = models.DateField()
    password = models.CharField(max_length=128)  # plaintext yerine hash önerilir

    def __str__(self):
        return self.email
