/* eslint-disable @typescript-eslint/no-explicit-any */
import elliptic from 'elliptic';
import crypto from 'crypto';

const SYSTEM_KEY = process.env.NEXT_PUBLIC_SYSTEM_KEY as string;

const WALLET_PADDING = process.env.NEXT_PUBLIC_WALLET_PADDING as string;

const ec = new elliptic.ec('p256');

export const EncryptPrivateKeyWithPassword = (
  privateKey: string,
  passsword: string,
): { cipherText: string; iv: string; authTag: string } => {
  const iv = crypto.randomBytes(12);
  const key = crypto.createHash('sha256').update(passsword).digest();
  const cipher = crypto.createCipheriv('aes-256-gcm', key, iv);

  const encrypted = Buffer.concat([
    cipher.update(privateKey, 'utf8'),
    cipher.final(),
  ]);
  const authTag = cipher.getAuthTag();

  return {
    cipherText: encrypted.toString('hex'),
    authTag: authTag.toString('hex'),
    iv: iv.toString('hex'),
  };
};

export const DecryptPrivateKeyWithPassword = (
  password: string,
  cipherText: string,
  ivHex: string,
  authTagHex: string,
): string | null => {
  try {
    const key = crypto.createHash('sha256').update(password).digest();
    const iv = Buffer.from(ivHex, 'hex');
    const authTag = Buffer.from(authTagHex, 'hex');
    const decipher = crypto.createDecipheriv('aes-256-gcm', key, iv);
    decipher.setAuthTag(authTag);

    const decrypted = Buffer.concat([
      decipher.update(Buffer.from(cipherText, 'hex')),
      decipher.final(),
    ]);

    return decrypted.toString('utf8');
  } catch {
    return null;
  }
};

export const EncryptPrivateKeyForExport = (privateKey: string): string => {
  const padded = `${WALLET_PADDING}:::${privateKey}`;

  const hmac = crypto.createHmac('sha256', SYSTEM_KEY).update(padded).digest();

  const iv = crypto.randomBytes(12);
  const key = crypto.createHash('sha256').update(SYSTEM_KEY).digest();
  const cipher = crypto.createCipheriv('aes-256-gcm', key, iv);
  const encryptedBuff = Buffer.concat([
    cipher.update(padded, 'utf8'),
    cipher.final(),
  ]);
  const authTag = cipher.getAuthTag();

  const encrypted = Buffer.concat([iv, authTag, encryptedBuff]).toString('hex');

  return `${encrypted}:::${hmac.toString('hex')}`;
};

export const DecryptedPrivateKeyFromExport = (
  encryptedHex: string,
): string | null => {
  try {
    const [cipherText, hmac] = encryptedHex.split(':::');
    const buf = Buffer.from(cipherText, 'hex');
    if (buf.length < 12 + 16) return null;

    const iv = buf.subarray(0, 12);
    const authTag = buf.subarray(12, 28);
    const cipherBuf = buf.subarray(28);

    const key = crypto.createHash('sha256').update(SYSTEM_KEY).digest();
    const decipher = crypto.createDecipheriv('aes-256-gcm', key, iv);
    decipher.setAuthTag(authTag);

    const decrypted = Buffer.concat([
      decipher.update(cipherBuf),
      decipher.final(),
    ]);
    const decoded = decrypted.toString('utf8');

    if (!decoded.startsWith(`${WALLET_PADDING}:::`)) return null;
    const privateKey = decoded.slice(`${WALLET_PADDING}:::`.length);

    const expectedHmac = crypto
      .createHmac('sha256', SYSTEM_KEY)
      .update(decoded)
      .digest('hex');

    if (hmac !== expectedHmac) return null;

    return privateKey;
  } catch {
    return null;
  }
};

export const SignPayload = (privateHex: string, data: any): string => {
  let dataStr = null;
  if (typeof data === 'string') {
    dataStr = data;
  } else {
    dataStr = JSON.stringify(data);
  }

  const hashBuffer = crypto.createHash('sha256').update(dataStr).digest();
  const key = ec.keyFromPrivate(privateHex, 'hex');
  const sig = key.sign(hashBuffer);

  const r = sig.r.toArray('be', 32);
  const s = sig.s.toArray('be', 32);

  return Buffer.concat([Buffer.from(r), Buffer.from(s)]).toString('hex');
};

export const SignTransaction = (
  privateHex: string,
  dataHex: string,
): string => {
  const dataToSign = Buffer.from(dataHex, 'hex');

  const key = ec.keyFromPrivate(privateHex, 'hex');
  const sig = key.sign(dataToSign);

  const r = sig.r.toArray('be', 32);
  const s = sig.s.toArray('be', 32);

  return Buffer.concat([Buffer.from(r), Buffer.from(s)]).toString('hex');
};
