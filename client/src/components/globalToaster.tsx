import { ReactNode } from 'react';
import { Action, Toaster, toast as baseToast } from 'sonner';

const ToastContent = ({
  icon,
  title,
  description,
}: {
  icon?: ReactNode;
  title: string;
  description?: string;
}) => (
  <div className="flex items-start gap-3">
    {icon && <div className="text-xl">{icon}</div>}
    <div className="flex flex-col">
      <p className="text-sm font-medium text-gray-900 dark:text-white">
        {title}
      </p>
      {description && (
        <p className="text-sm text-gray-600 dark:text-gray-400">
          {description}
        </p>
      )}
    </div>
  </div>
);

const baseStyle =
  'flex items-center gap-3 px-5 py-4 bg-white dark:bg-primary-dark rounded-xl shadow-lg transition-colors backdrop-blur-sm';

export const toast = {
  success: (
    title: string,
    desc?: string,
    icon?: ReactNode,
    action?: ReactNode | Action,
  ) =>
    baseToast.custom(
      () => (
        <ToastContent icon={icon ?? '✅'} title={title} description={desc} />
      ),
      {
        action,
        className: `${baseStyle} text-white`,
      },
    ),

  error: (
    title: string,
    desc?: string,
    icon?: ReactNode,
    action?: ReactNode | Action,
  ) =>
    baseToast.custom(
      () => (
        <ToastContent icon={icon ?? '❌'} title={title} description={desc} />
      ),
      {
        action,
        className: `${baseStyle} text-white`,
      },
    ),

  warning: (
    title: string,
    desc?: string,
    icon?: ReactNode,
    action?: ReactNode | Action,
  ) =>
    baseToast.custom(
      () => (
        <ToastContent icon={icon ?? '⚠️'} title={title} description={desc} />
      ),
      {
        action,
        className: `${baseStyle} text-white`,
      },
    ),

  info: (
    title: string,
    desc?: string,
    icon?: ReactNode,
    action?: ReactNode | Action,
  ) =>
    baseToast.custom(
      () => (
        <ToastContent icon={icon ?? 'ℹ️'} title={title} description={desc} />
      ),
      {
        action,
        className: `${baseStyle} text-white`,
      },
    ),
};

export const GlobalToaster = () => (
  <Toaster position="top-right" toastOptions={{ duration: 4000 }} />
);
