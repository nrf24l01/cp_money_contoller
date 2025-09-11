from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from time import sleep
import os


def register_user(driver: webdriver.Chrome, invite: str, email: str, password: str, name: str, birthday: str, learn_place: str, grade: int):
    """Регистрирует пользователя на сайте с автоматическим решением капчи"""
    driver.get("https://codingprojects.ru/register")
    
    # Заполняем форму
    driver.find_element("name", "invite").send_keys(invite)
    
    driver.find_element("name", "email").send_keys(email)
    driver.find_element("name", "password").send_keys(password)
    driver.find_element("name", "password_confirmation").send_keys(password)
    
    driver.find_element("name", "name").send_keys(name)
    driver.find_element("name", "birthday").send_keys(birthday)
    driver.find_element("name", "school").send_keys(learn_place)
    driver.find_element("name", "grade").send_keys(str(grade))
    
    # Ждем загрузки рекапчи и переключаемся на iframe
    wait = WebDriverWait(driver, 10)
    iframe = wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, "iframe[title='reCAPTCHA']")))
    driver.switch_to.frame(iframe)
    
    # Кликаем на чекбокс рекапчи
    recaptcha_checkbox = wait.until(EC.element_to_be_clickable((By.CLASS_NAME, "recaptcha-checkbox-border")))
    recaptcha_checkbox.click()
    
    # Переключаемся обратно к основному контенту
    driver.switch_to.default_content()
    
    # Ждем решения капчи расширением (максимум 30 секунд)
    print("Ожидание решения капчи расширением...")
    for i in range(30):
        try:
            # Проверяем, появилась ли галочка в рекапче
            driver.switch_to.frame(iframe)
            checkbox = driver.find_element(By.CLASS_NAME, "recaptcha-checkbox-checkmark")
            if checkbox.is_displayed():
                print("Капча решена!")
                driver.switch_to.default_content()
                break
            driver.switch_to.default_content()
        except:
            pass
        
        sleep(1)
        print(f"Ожидание... {i+1}/30 секунд")
    else:
        print("Капча не была решена автоматически")
        driver.switch_to.default_content()
    
    # Пытаемся нажать кнопку регистрации
    try:
        submit_button = driver.find_element(By.CSS_SELECTOR, "button[type='submit'], input[type='submit']")
        submit_button.click()
        print("Форма отправлена!")
    except Exception as e:
        print(f"Не удалось найти кнопку отправки: {e}")
    
    # Ждем результата
    sleep(5)


if __name__ == "__main__":
    from selenium import webdriver
    from selenium.webdriver.chrome.service import Service
    from webdriver_manager.chrome import ChromeDriverManager
    from . import build_config
    
    # Настраиваем Chrome с расширением
    options, service = build_config()
    
    driver = webdriver.Chrome(service=service, options=options)
    register_user(driver, "invite_code", "fake_email@example.com", "fake_password", "Fake Name", "1990-01-01", "Fake School", 10)

