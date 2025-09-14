from imapclient import IMAPClient
import email

class ImapClient:
    def __init__(self, host: str, username: str, password: str, ssl: bool = True):
        self.client = IMAPClient(host, use_uid=True, ssl=ssl)
        self.client.login(username, password)

    def find_from(self, from_address: str):
        self.client.select_folder('INBOX')
        messages = self.client.search(['FROM', from_address, 'UNSEEN'])
        if not messages:
            return None
        # Get the latest unread message by UID (assuming UIDs are sequential)
        latest_uid = max(messages)
        raw_message = self.client.fetch([latest_uid], ['RFC822'])[latest_uid][b'RFC822']
        email_message = email.message_from_bytes(raw_message)
        body = None
        if email_message.is_multipart():
            for part in email_message.walk():
                if part.get_content_type() == "text/plain":
                    body = part.get_payload(decode=True).decode('utf-8')
                    break
        else:
            body = email_message.get_payload(decode=True).decode('utf-8')
        # Mark the message as read
        self.client.add_flags([latest_uid], ['SEEN'])
        return body

    def logout(self):
        self.client.logout()