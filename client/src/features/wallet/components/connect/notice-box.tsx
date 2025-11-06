type NoticeVariant = 'info' | 'warning' | 'success' | 'error';

interface NoticeBoxProps {
  variant: NoticeVariant;
  title: string;
  description: string;
  icon: React.ReactNode;
  style?: React.CSSProperties;
}

const NoticeBox = ({
  description,
  icon,
  title,
  variant,
  style,
}: NoticeBoxProps) => {
  const color = {
    bgColor: '',
    borderColor: '',
    titleColor: '',
    desColor: '',
  };

  switch (variant) {
    case 'success':
      color.bgColor = 'bg-green-50';
      color.borderColor = 'border-green-600';
      color.titleColor = 'text-green-600';
      color.desColor = 'text-green-500';
      break;
    case 'warning':
      break;
    case 'info':
      break;
    case 'error':
      color.bgColor = 'bg-red-50';
      color.borderColor = 'border-red-600';
      color.titleColor = 'text-red-600';
      color.desColor = 'text-red-500';
      break;
  }

  return (
    <div
      style={style}
      className={`animate-cascase-fade glass-card ${color.bgColor} rounded-xl p-4 border-l-4 ${color.borderColor}`}
    >
      <div className="flex items-start space-x-3">
        {icon}
        <div>
          <p className={`${color.titleColor} font-semibold text-sm`}>{title}</p>
          <p
            className={`${color.desColor} opacity-80 text-xs mt-1 leading-relaxed`}
          >
            {description}
          </p>
        </div>
      </div>
    </div>
  );
};

export default NoticeBox;
