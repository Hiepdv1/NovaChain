import Button from '@/components/button';
import { toast } from '@/components/globalToaster';
import Input from '@/components/input';
import {
  ChangeEvent,
  Fragment,
  memo,
  useCallback,
  useEffect,
  useRef,
} from 'react';
import NoticeBox from './notice-box';

type RefMap = {
  strengthDiv: HTMLDivElement | null;
  strengthText: HTMLSpanElement | null;
  strengthFill: HTMLDivElement | null;
  matchText: HTMLSpanElement | null;
  matchDiv: HTMLDivElement | null;
  createPassword: HTMLInputElement | null;
  confirmPassword: HTMLInputElement | null;
  createPasswordBtn: HTMLButtonElement | null;
};

type StrengthStatus = {
  level: 'weak' | 'fair' | 'good' | 'strong';
  text: string;
  class: string;
};

interface CreatePasswordProps {
  onBack: () => void;
  onContinue: (password: string) => void;
  preview?: React.ReactNode;
}

const CreatePassword = ({
  onBack,
  onContinue,
  preview,
}: CreatePasswordProps) => {
  const refs = useRef<RefMap>({
    matchText: null,
    strengthDiv: null,
    strengthText: null,
    strengthFill: null,
    matchDiv: null,
    createPassword: null,
    confirmPassword: null,
    createPasswordBtn: null,
  });

  const checkPasswordStrength = (password: string): StrengthStatus => {
    let strength = 0;
    if (password.length >= 8) strength++;
    if (password.length >= 12) strength++;
    if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++;
    if (/\d/.test(password)) strength++;
    if (/[^a-zA-Z0-9]/.test(password)) strength++;

    switch (strength) {
      case 0:
      case 1:
        return { class: 'strength-weak', level: 'weak', text: 'Weak' };
      case 2:
        return { class: 'strength-fair', level: 'fair', text: 'Fair' };
      case 3:
      case 4:
        return { class: 'strength-good', level: 'good', text: 'Good' };
      case 5:
        return { class: 'strength-strong', level: 'strong', text: 'Strong' };
      default:
        return { class: 'strength-week', level: 'weak', text: 'Week' };
    }
  };

  const checkCreatePasswordMatch = useCallback(() => {
    const {
      createPassword,
      confirmPassword,
      matchDiv,
      matchText,
      createPasswordBtn,
    } = refs.current;

    const password = createPassword?.value || '';
    const confirmPasswordValue = confirmPassword?.value || '';

    if (confirmPasswordValue.length > 0) {
      matchDiv?.classList.remove('hidden');

      if (password === confirmPasswordValue) {
        if (matchText) {
          matchText.textContent = '✓ Passwords match';
          matchText.className = 'font-bold text-sm text-emerald-600';
        }

        if (
          checkPasswordStrength(password).level !== 'weak' &&
          createPasswordBtn
        ) {
          createPasswordBtn.disabled = false;
        }
      } else {
        if (matchText) {
          matchText.textContent = '✗ Passwords do not match';
          matchText.className = 'font-bold text-sm text-red-600';
        }
        if (createPasswordBtn) {
          createPasswordBtn.disabled = true;
        }
      }
    } else {
      matchDiv?.classList.add('hidden');
      if (createPasswordBtn) {
        createPasswordBtn.disabled = true;
      }
    }
  }, []);

  const onProceedToWallet = useCallback(() => {
    const { createPassword, confirmPassword } = refs.current;

    const password = createPassword?.value || '';
    const confirmPasswordValue = confirmPassword?.value || '';

    if (password !== confirmPasswordValue) {
      toast.error('Passwords do not match');
      return;
    }

    if (checkPasswordStrength(password).level === 'weak') {
      toast.error('Please use a stronger password');
      return;
    }

    onContinue(password);
  }, [onContinue]);

  const onInputCreatePassword = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      const password = e.target.value;
      const { strengthFill, strengthText, strengthDiv } = refs.current;

      if (password.length > 0) {
        strengthDiv?.classList.remove('hidden');

        const strength = checkPasswordStrength(password);

        if (strengthFill) {
          strengthFill.className = `strength-bar-fill ${strength.class}`;
        }

        if (strengthText) {
          strengthText.textContent = strength.text;

          const colorClass =
            strength.level === 'strong'
              ? 'text-emerald-600'
              : strength.level === 'good'
              ? 'text-blue-600'
              : strength.level === 'fair'
              ? 'text-yellow-600'
              : 'text-red-600';

          strengthText.className = `font-bold ${colorClass}`;
        }
      } else {
        strengthDiv?.classList.add('hidden');
      }

      checkCreatePasswordMatch();
    },
    [checkCreatePasswordMatch],
  );

  useEffect(() => {
    const { createPasswordBtn } = refs.current;
    if (createPasswordBtn) {
      createPasswordBtn.disabled = true;
    }
  }, []);

  return (
    <Fragment>
      <div className="space-y-10">
        <div className="text-center mb-6 animate-cascase-fade delay-200">
          <h2 className="text-3xl font-black bg-gradient-to-r from-slate-300 via-white to-slate-100 bg-clip-text text-transparent mb-6">
            Create Your Password
          </h2>
          <p className="text-slate-200 text-sm font-medium">
            This password keeps your wallet secure on this device. Only you can
            unlock it - but you can reset it anytime with your private key.
          </p>
        </div>

        {preview}

        <div className="space-y-5">
          <div
            style={{
              animationDelay: '300ms',
            }}
            className="enhanced-floating animate-cascase-fade"
          >
            <Input
              id="createPassword"
              variant="levitating"
              inputSize="sm"
              placeholder=""
              type="password"
              ref={(el) => void (refs.current.createPassword = el)}
              onInput={onInputCreatePassword}
            />
            <label className="text-sm" htmlFor="createPassword">
              Enter Password
            </label>
          </div>

          <div
            ref={(el) => void (refs.current.strengthDiv = el)}
            style={{
              animationDelay: '400ms',
            }}
            className="animate-cascase-fade text-xs hidden"
          >
            <div className="flex justify-between items-center mb-3">
              <span className="text-slate-800 font-semibold">
                Password Strength
              </span>
              <span
                ref={(el) => void (refs.current.strengthText = el)}
                className="font-bold"
              >
                Weak
              </span>
            </div>
            <div className="strength-evolution">
              <div
                ref={(el) => void (refs.current.strengthFill = el)}
                className="strength-bar-fill strength-weak"
              ></div>
            </div>
          </div>

          <div
            className="enhanced-floating animate-cascase-fade"
            style={{
              animationDelay: '500ms',
            }}
          >
            <Input
              id="confirmPassword"
              variant="levitating"
              inputSize="sm"
              placeholder=""
              type="password"
              ref={(el) => void (refs.current.confirmPassword = el)}
              onInput={checkCreatePasswordMatch}
            ></Input>
            <label className="text-sm" htmlFor="confirmPassword">
              Confirm Password
            </label>
          </div>

          <div
            ref={(el) => void (refs.current.matchDiv = el)}
            className="hidden text-center"
          >
            <span
              ref={(el) => void (refs.current.matchText = el)}
              className="font-bold text-sm text-red-600"
            ></span>
          </div>
        </div>

        <NoticeBox
          description="Your password should be strong: include uppercase, lowercase, numbers, and special characters for maximum security."
          icon={
            <svg
              className="w-10 h-10 text-green-600 mt-1"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M12,1L3,5V11C3,16.55 6.84,21.74 12,23C17.16,21.74 21,16.55 21,11V5L12,1M9,12L7,10L5.5,11.5L9,15L18.5,5.5L17,4L9,12Z"></path>
            </svg>
          }
          title="🔐 Strong Password Required"
          variant="success"
          style={{
            animationDelay: '600ms',
          }}
        />

        <NoticeBox
          description="Your private key will be encrypted with this password and stored
                securely in your browser. Only you will have access to decrypt
                it."
          icon={
            <svg
              className="w-10 h-10 text-green-600 mt-1"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M12,1L3,5V11C3,16.55 6.84,21.74 12,23C17.16,21.74 21,16.55 21,11V5L12,1M9,12L7,10L5.5,11.5L9,15L18.5,5.5L17,4L9,12Z"></path>
            </svg>
          }
          title="🔐 Encryption Notice"
          variant="success"
          style={{
            animationDelay: '600ms',
          }}
        />

        <div
          className="flex space-x-6 animate-cascase-fade"
          style={{
            animationDelay: '700ms',
          }}
        >
          <Button onClick={onBack} variant="glass" size="md">
            ← Back
          </Button>
          <Button
            onClick={onProceedToWallet}
            ref={(el) => void (refs.current.createPasswordBtn = el)}
            variant="quantum"
            size="md"
          >
            Continue →
          </Button>
        </div>
      </div>
    </Fragment>
  );
};

export default memo(CreatePassword);
