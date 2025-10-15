import crypto from 'crypto';

const SECRET_KEY = process.env.NEXT_PUBLIC_ENCODE_DATA_SECRET_KEY as string;

function getKey() {
  return crypto.createHash('sha256').update(SECRET_KEY).digest();
}

export function EncryptData<T>(data: T): string {
  const json = JSON.stringify(data);
  const key = getKey();

  const iv = crypto.randomBytes(12);

  const cipher = crypto.createCipheriv('aes-256-gcm', key, iv);

  const encrypted = Buffer.concat([
    cipher.update(json, 'utf8'),
    cipher.final(),
  ]);

  const tag = cipher.getAuthTag();

  const buffer = Buffer.concat([iv, tag, encrypted]);

  return buffer.toString('base64');
}

export function DecryptData<T>(encrypted: string): T {
  const raw = Buffer.from(encrypted, 'base64');

  const key = getKey();
  const iv = raw.subarray(0, 12);
  const tag = raw.subarray(12, 28);
  const ciphertext = raw.subarray(28);

  const decipher = crypto.createDecipheriv('aes-256-gcm', key, iv);
  decipher.setAuthTag(tag);

  const decrypted = Buffer.concat([
    decipher.update(ciphertext),
    decipher.final(),
  ]);

  return JSON.parse(decrypted.toString('utf8'));
}
