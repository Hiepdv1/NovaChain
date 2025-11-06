import { memo } from 'react';

export interface LoadingStep {
  active: boolean;
  completed: boolean;
  title: string;
}

export interface LoadingOverlayProps {
  title: string;
  des: string;
  loadingSteps: LoadingStep[];
}

const LoadingOverlay = ({ des, title, loadingSteps }: LoadingOverlayProps) => {
  return (
    <div className="loading-overlay">
      <div className="loading-content">
        <h3 className="text-xl font-bold text-slate-800 mb-2">{title}</h3>
        <p className="text-slate-600 mb-4 text-sm">{des}</p>
        <div className="loading-dots">
          <div
            style={{
              animationDelay: '-0.32s',
            }}
            className="loading-dot"
          ></div>
          <div
            style={{
              animationDelay: '-0.16s',
            }}
            className="loading-dot"
          ></div>
          <div
            style={{
              animationDelay: '0s',
            }}
            className="loading-dot"
          ></div>
        </div>

        <div className="loading-steps">
          {loadingSteps.map((step, ix) => {
            return (
              <div
                key={ix}
                className={`loading-step ${step.completed && 'completed'} ${
                  step.active && 'active'
                }`}
              >
                <div className="loading-step-number">{ix + 1}</div>
                <div className="loading-step-text">{step.title}</div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};

export default memo(LoadingOverlay);
