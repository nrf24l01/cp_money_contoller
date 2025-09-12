from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from time import sleep

def register_user(driver: webdriver.Chrome, invite: str, email: str, password: str, name: str, birthday: str, learn_place: str, grade: int):
    """Регистрирует пользователя на сайте с автоматическим решением капчи"""
    driver.get("https://codingprojects.ru/register")
    # Скриншот всей страницы

    # Заполняем форму
    driver.find_element("name", "invite").send_keys(invite)
    
    driver.find_element("name", "email").send_keys(email)
    driver.find_element("name", "password").send_keys(password)
    driver.find_element("name", "password_confirmation").send_keys(password)
    
    driver.find_element("name", "name").send_keys(name)
    driver.find_element("name", "birthday").send_keys(birthday)
    driver.find_element("name", "school").send_keys(learn_place)
    driver.find_element("name", "grade").send_keys(str(grade))
    
    for j in range(3):
        response = driver.execute_script("return document.getElementById('g-recaptcha-response').value")
        
        if response:
            break
        print(f"Попытка решения капчи {j+1}/3")
        wait = WebDriverWait(driver, 10)
        iframe = wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, "iframe[title='reCAPTCHA']")))
        driver.switch_to.frame(iframe)
        recaptcha_checkbox = wait.until(EC.element_to_be_clickable((By.CLASS_NAME, "recaptcha-checkbox-border")))
        recaptcha_checkbox.click()
        driver.switch_to.default_content()
        
        print("Ожидание решения капчи...")
        for i in range(30):
            response = driver.execute_script("return document.getElementById('g-recaptcha-response').value")
            if response:
                print("CAPTCHA решена")
                break
            sleep(1)
            print(f"Ожидание... {i+1}/30 секунд")
            driver.save_screenshot("waiting.png")
            try:
                challenge_iframe = driver.find_element(By.CSS_SELECTOR, "iframe[title='recaptcha challenge expires in two minutes']")
                driver.switch_to.frame(challenge_iframe)
                try:
                    reset_button = driver.find_element(By.ID, "reset-button")
                    reset_button.click()
                    print("Clicked reset button")
                    break
                except:
                    pass
                html = driver.page_source
                with open("recaptcha_challenge.html", "w", encoding="utf-8") as f:
                    f.write(html)
                print("HTML saved to recaptcha_challenge.html")
                help_button = WebDriverWait(driver, 1).until(EC.element_to_be_clickable((By.CSS_SELECTOR, ".button-holder.help-button-holder")))
                help_button.click()
                print("Clicked autosolve button")
            except:
                pass
            finally:
                driver.switch_to.default_content()
        
        
    else:
        print("Капча не была решена автоматически")
    
    # Пытаемся нажать кнопку регистрации
    try:
        submit_button = driver.find_element(By.CSS_SELECTOR, "button[type='submit'], input[type='submit']")
        submit_button.click()
        print("Форма отправлена!")
    except Exception as e:
        print(f"Не удалось найти кнопку отправки: {e}")
    
    # Проверяем наличие ошибки
    try:
        error_div = driver.find_element(By.CSS_SELECTOR, "div.alert.alert-danger.alert-dismissible")
        error_text = error_div.text
        print(error_text[2:])
    except:
        print("Регистрация прошла успешно или ошибка не найдена.")
    
    try:
        driver.find_element(By.XPATH, "/html/body/div[1]/div[2]/div/div/div/div/div/div/div")
        driver.save_screenshot("complete.png")
        print("Регистрация удалась")
    except:
        print("Регистрация не удалась")
        # Screenshot for debugging
        driver.save_screenshot("registration_error.png")
        print("Скриншот сохранен как registration_error.png")
    
    # Ждем результата
    sleep(100)


if __name__ == "__main__":
    from .preconfigure import build_config
    
    # Настраиваем Chrome с расширением
    options, service = build_config(headless=True)
    
    driver = webdriver.Chrome(service=service, options=options)
    register_user(driver, "RSaEMYwq", "fake_email+fak@example.com", "fake_password", "Fake Name", "1990-01-01", "Fake School", 10)

