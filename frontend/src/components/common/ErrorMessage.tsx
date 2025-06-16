import React from 'react';
import { FaExclamationCircle } from 'react-icons/fa';

interface ErrorMessageProps {
  message: string;
  className?: string;
  onRetry?: () => void;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({
  message,
  className = '',
  onRetry
}) => {
  return (
    <div
      className={`bg-red-50 text-red-700 flex flex-col items-center justify-center rounded-lg p-4 ${className}`}
    >
      <div className="flex items-center">
        <FaExclamationCircle className="mr-2 text-xl" />
        <span className="font-medium">{message}</span>
      </div>
      {onRetry && (
        <button
          onClick={onRetry}
          className="bg-red-600 hover:bg-red-700 mt-2 rounded px-4 py-1 text-sm text-white transition-colors"
        >
          Retry
        </button>
      )}
    </div>
  );
};

export default ErrorMessage;
