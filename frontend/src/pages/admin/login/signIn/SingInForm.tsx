import { useState } from 'react';
import { useForm, SubmitHandler } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { singInSchema } from './schema/sindInSchema';
import { FormValuesSignIn } from '@/types';
import { useAuth } from '@/contexts/AuthContext';

interface ApiError extends Error {
  response?: {
    status: number;
    data?: {
      message?: string;
    };
  };
}

const SingInForm = () => {
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();

  const {
    register,
    handleSubmit,
    setError,
    formState: { errors, isValid }
  } = useForm<FormValuesSignIn>({
    resolver: zodResolver(singInSchema),
    mode: 'onChange'
  });

  const onSubmit: SubmitHandler<FormValuesSignIn> = async (data) => {
    setIsLoading(true);
    try {
      await login({
        email: data.email,
        password: data.password
      });
      // Navigation happens automatically after successful login in the AuthProvider
    } catch (error) {
      const apiError = error as ApiError;
      const errorMessage =
        apiError?.response?.data?.message || 'An error occurred during login';

      if (apiError.response?.status === 404) {
        setError('email', {
          type: 'manual',
          message: errorMessage
        });
      } else {
        setError('password', {
          type: 'manual',
          message: errorMessage
        });
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleForgotPassword = () => {
    // TODO: Implement forgot password functionality
    console.log('Forgot password clicked');
  };

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className="mx-auto mb-[60px] w-[330px]"
    >
      <div className="mb-4">
        <label htmlFor="email" className="mb-2 block">
          Email:
        </label>
        <input
          id="email"
          type="email"
          {...register('email')}
          className={`w-full border ${
            errors.email ? 'border-red' : 'border-darkgray'
          } rounded px-4 py-2`}
          placeholder="Введіть свій email"
        />
        {errors.email && (
          <p className="text-red-500 mt-1 text-sm">{errors.email.message}</p>
        )}
      </div>

      <div className="mb-6">
        <label htmlFor="password" className="mb-2 block">
          Пароль:
        </label>
        <input
          id="password"
          type="password"
          {...register('password')}
          className={`w-full border ${
            errors.password ? 'border-red' : 'border-darkgray'
          } rounded px-4 py-2`}
          placeholder="Введіть пароль"
        />
        {errors.password && (
          <p className="text-red-500 mt-1 text-sm">{errors.password.message}</p>
        )}
      </div>

      <button
        type="button"
        onClick={handleForgotPassword}
        className="mb-6 block text-blue-500 hover:underline"
      >
        Не пам'ятаю пароль
      </button>

      <button
        type="submit"
        disabled={!isValid || isLoading}
        className={`w-full rounded bg-blue-500 px-4 py-2 text-white transition-colors hover:bg-blue-600 ${
          !isValid || isLoading ? 'cursor-not-allowed opacity-50' : ''
        }`}
      >
        {isLoading ? 'Вхід...' : 'Увійти'}
      </button>
    </form>
  );
};

export default SingInForm;
