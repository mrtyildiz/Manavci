from django.shortcuts import render, redirect
from django.contrib import messages
# Create your views here.
from .models import ManavciUser


def login(request):
    if request.method == 'POST':
        email = request.POST.get('email')
        password = request.POST.get('password')

        try:
            user = ManavciUser.objects.get(email=email)
            if user.password == password:  # ⚠️ hashlenmiş şifre kontrolü önerilir
                request.session['user_id'] = user.id
                request.session['user_email'] = user.email
                return redirect('index')  # login sonrası yönlendir
            else:
                messages.error(request, "Invalid password.")
        except ManavciUser.DoesNotExist:
            messages.error(request, "User not found.")
    return render(request, 'login.html')
    # return render(request, 'login.html')

def index(request):
    if not request.session.get('user_id'):
        return redirect('login')
    return render(request, 'index.html', {'user_email': request.session.get('user_email')})

def register(request):
    context = {}
    if request.method == 'POST':
        data = request.POST
        first_name = data.get('firstName')
        last_name = data.get('lastName')
        email = data.get('email')
        gender = data.get('gender')
        dob = data.get('dob')
        password = data.get('password')
        confirm_password = data.get('confirmPassword')

        context['confirm_password'] = confirm_password  # Geri gösterim için

        if password != confirm_password:
            messages.error(request, "Passwords do not match.")
            print("Şifreler uyuşmuyor")
            return render(request, 'register.html', context)

        if ManavciUser.objects.filter(email=email).exists():
            messages.error(request, "Email already exists.")
            return render(request, 'register.html', context)

        ManavciUser.objects.create(
            first_name=first_name,
            last_name=last_name,
            email=email,
            gender=gender,
            date_of_birth=dob,
            password=password  # ⚠️ Hash önerilir
        )

        messages.success(request, "User registered successfully.")
        return redirect('login')

    return render(request, 'register.html', context)


def send_mail(request):
    return render(request, 'SendMail.html')

def logout_view(request):
    request.session.flush()
    messages.success(request, "Logged out.")
    return redirect('login')