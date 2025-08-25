import crypto from 'crypto';
import { ec as EC } from 'elliptic';
import RIPEMD160 from 'ripemd160';
import bs58 from 'bs58';

export interface KeyPair {
  privateKey: string;
  publicKey: string;
}

const version = 0x00;
const checkSumLength = 4;

export const CreateWallet = (): KeyPair => {
  const ec = new EC('p256');
  const keyPair = ec.genKeyPair();

  const privateKey = keyPair.getPrivate().toArray();
  const publicKey = Buffer.concat([
    Buffer.from(keyPair.getPublic().getX().toArray()),
    Buffer.from(keyPair.getPublic().getY().toArray()),
  ]);

  return {
    privateKey: Buffer.from(privateKey).toString('hex'),
    publicKey: publicKey.toString('hex'),
  };
};

export const GetPublicKeyFromPrivateKey = (privateKey: string): string => {
  const ec = new EC('p256');
  const privKey = Buffer.from(privateKey, 'hex');
  const keyPair = ec.keyFromPrivate(privKey);
  const publicKey = Buffer.concat([
    Buffer.from(keyPair.getPublic().getX().toArray()),
    Buffer.from(keyPair.getPublic().getY().toArray()),
  ]);

  return publicKey.toString('hex');
};

export const PublicKeyHash = (publicKey: string): Buffer => {
  const pubKeyBytes = Buffer.from(publicKey, 'hex');
  const sha256 = crypto.createHash('sha256').update(pubKeyBytes).digest();

  return new RIPEMD160().update(sha256).digest();
};

const createCheckSum = (data: Buffer): Buffer => {
  const firstHash = crypto.createHash('sha256').update(data).digest();
  const secondHash = crypto.createHash('sha256').update(firstHash).digest();
  return secondHash.subarray(0, checkSumLength);
};

export const GetAddress = (publicKey: string): string => {
  const pubHash = PublicKeyHash(publicKey);
  const versionedHash = Buffer.concat([Buffer.from([version]), pubHash]);
  const checksum = createCheckSum(versionedHash);
  const fullHash = Buffer.concat([versionedHash, checksum]);
  return bs58.encode(fullHash);
};
