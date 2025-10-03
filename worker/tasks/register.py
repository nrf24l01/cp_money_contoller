import re
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from time import sleep, time
from modules import ThreadSafeLogger, ImapClient
from threading import Thread
import traceback


def create_user(driver: webdriver.Chrome, logger: ThreadSafeLogger, invite: str, email: str, password: str, name: str, birthday: str, learn_place: str, grade: int, resp: dict):
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
    
    try:
        for j in range(3):
            response = driver.execute_script("return document.getElementById('g-recaptcha-response').value")
            
            if response:
                break
            logger.info(f"Попытка решения капчи {j+1}/3")
            wait = WebDriverWait(driver, 10)
            iframe = wait.until(EC.presence_of_element_located((By.CSS_SELECTOR, "iframe[title='reCAPTCHA']")))
            driver.switch_to.frame(iframe)
            recaptcha_checkbox = wait.until(EC.element_to_be_clickable((By.CLASS_NAME, "recaptcha-checkbox-border")))
            recaptcha_checkbox.click()
            driver.switch_to.default_content()
            logger.info("Ожидание решения капчи...")
            for i in range(30):
                response = driver.execute_script("return document.getElementById('g-recaptcha-response').value")
                if response:
                    logger.info("CAPTCHA решена")
                    break
                sleep(1)
                logger.info(f"Ожидание... {i+1}/30 секунд")
                driver.save_screenshot("waiting.png")
                try:
                    challenge_iframe = driver.find_element(By.CSS_SELECTOR, "iframe[title='recaptcha challenge expires in two minutes']")
                    driver.switch_to.frame(challenge_iframe)
                    try:
                        reset_button = driver.find_element(By.ID, "reset-button")
                        reset_button.click()
                        logger.info("Clicked reset button")
                        break
                    except:
                        pass
                    help_button = WebDriverWait(driver, 1).until(EC.element_to_be_clickable((By.CSS_SELECTOR, ".button-holder.help-button-holder")))
                    help_button.click()
                    logger.info("Clicked autosolve button")
                except:
                    pass
                finally:
                    driver.switch_to.default_content()
    except Exception as e:
        logger.error(f"Ошибка при решении капчи: {traceback.format_exc()}")
        resp["res"] = False
    # Пытаемся нажать кнопку регистрации
    try:
        submit_button = driver.find_element(By.CSS_SELECTOR, "button[type='submit'], input[type='submit']")
        submit_button.click()
        logger.info("Форма отправлена!")
    except Exception as e:
        logger.error(f"Не удалось найти кнопку отправки: {e}")
    # Проверяем наличие ошибки
    try:
        error_div = driver.find_element(By.CSS_SELECTOR, "div.alert.alert-danger.alert-dismissible")
        error_text = error_div.text
        logger.error(error_text[2:])
    except:
        logger.info("Регистрация прошла успешно или ошибка не найдена.")
        
    try:
        submit_button = driver.find_element(By.CSS_SELECTOR, "button[type='submit'], input[type='submit']")
        submit_button.click()
        logger.info("Форма отправлена!")
    except Exception as e:
        logger.warn(f"Не удалось найти кнопку отправки: {e}")
        
    try:
        driver.find_element(By.XPATH, "/html/body/div[1]/div[2]/div/div/div/div/div/div/div")
        driver.save_screenshot("complete.png")
        logger.info("Регистрация удалась, необходимо подтвердить email")
        resp["res"] = True
        return True
    except:
        logger.error("Регистрация не удалась")
        # Screenshot for debugging
        driver.save_screenshot("registration_error.png")
        logger.info("Скриншот сохранен как registration_error.png")
        try:
            error_blocks = driver.find_elements(By.CLASS_NAME, "error-block")
            for error_block in error_blocks:
                logger.error(f"Error block: {error_block.text}")
        except Exception as e:
            logger.error(f"Failed to find error blocks: {e}")
        resp["res"] = False
        return False

def verify_email(driver: webdriver.Chrome, logger: ThreadSafeLogger, link: str, email: str, password: str, result: dict):
    logger.info(f"Переход по ссылке верификации: {link}")
    driver.get(link)
    
    try:
        driver.find_element(By.ID, "inputEmail").send_keys(email)
        driver.find_element(By.ID, "inputPassword").send_keys(password)
        driver.find_element(By.CSS_SELECTOR, "button[type='submit']").click()
    except:
        pass
    
    try:
        driver.find_element(By.XPATH, "//img[@src='/images/clip-education.png']")
        result['verified'] = True
        logger.info("Email успешно верифицирован")
        return True
    except:
        result['verified'] = False
        logger.error("Не удалось верифицировать email")
        return False

def wait_for_mail(logger: ThreadSafeLogger, imap_client: ImapClient, verify_link: dict, timeout: int = 300,  email: str = "noreply@codingprojects.ru"):
    logger.info(f"Waiting for email for {email}")
    start_time = time()
    while time() - start_time < timeout:
        email_message = imap_client.find_from(email)
        if email_message:
            pattern = r'https://codingprojects\.ru/email/verify/\d+/[a-f0-9]+(?:\?[^ \]\r\n]*)?'
            matches = re.findall(pattern, email_message)
            logger.info(f"Verification link found: {matches[0]}")
            logger.info("Email received")
            verify_link['url'] = matches[0]
            return matches[0]
        sleep(5)
    logger.warn("Email not received within timeout")
    return None

def register_user(driver: webdriver.Chrome, logger: ThreadSafeLogger, imap_client: ImapClient, invite: str, email: str, password: str, name: str, birthday: str, learn_place: str, grade: int, res: dict = {}):
    logger.info("Start registration process")

    resp = {"res": False}
    verify_link = {}
    email_thread = Thread(target=wait_for_mail, args=(logger, imap_client, verify_link), name="EmailThread", daemon=True)
    selenium_thread = Thread(target=create_user, args=(driver, logger, invite, email, password, name, birthday, learn_place, grade, resp), name="CreateUserThread", daemon=True)
    
    selenium_thread.start()
    email_thread.start()
    
    selenium_thread.join()
    if resp["res"] is False:
        logger.error("User creation failed, aborting verification")
        res["success"] = False
        return False
    email_thread.join()

    link = verify_link.get('url', None)
    if link is None:
        logger.error("Failed to get verification link from email")
        res["success"] = False
        return False
    result = {}
    
    verify_thread = Thread(target=verify_email, args=(driver, logger, link, email, password, result), name="VerifyEmailThread", daemon=True)
    verify_thread.start()
    
    logger.info("Waiting for email verification to complete")
    verify_thread.join()
    logger.info("Email verification process completed")
    logger.info("Waiting for email thread to finish")
    email_thread.join()
    logger.info("Email thread finished")
    
    if result.get('verified', False):
        logger.info("User created and verified successfully")
        res["success"] = True
        return True
    else:
        logger.error("User created but verification failed")
        res["success"] = False
        return False
    

if __name__ == "__main__":
    from .preconfigure import build_config
    import os
    
    logger = ThreadSafeLogger(__name__)
    
    options, service = build_config(headless=True)
    driver = webdriver.Chrome(service=service, options=options)

    IMAP_HOST = os.getenv("IMAP_HOST")
    IMAP_USER = os.getenv("IMAP_USER")
    IMAP_PASSWORD = os.getenv("IMAP_PASSWORD")

    imap_client = ImapClient(IMAP_HOST, IMAP_USER, IMAP_PASSWORD)

    register_user(driver, logger, imap_client, "",  "", "", "", "2000-01-01", "СШ №1", 10)

