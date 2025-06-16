/* eslint-disable prettier/prettier */
import TextInput from '@/components/admin/inputs/TextInput';
import { ContactsFormInput } from '@/types';
import { Dispatch, SetStateAction, useEffect, useState } from 'react';
import { Controller, SubmitHandler, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useAppDispatch } from '@/store/hook';
import { editEmail, editPhone, editName } from '@/store/slices/contactsSlice';
import { contactsValidation } from './contactsSchema';
import { TfiClose } from 'react-icons/tfi';
import { openAlert } from '@/store/slices/responseAlertSlice';
import {
  editErrorResponseMessage,
  editSuccessResponseMessage
} from '@/utils/responseMessages';
import { useNavigate } from 'react-router-dom';

type EditContactsProps = {
  id: string;
  data: string;
  currentType: 'name' | 'email' | 'phone';
  setIsModalOpen: Dispatch<SetStateAction<boolean>>;
};

const Edit = ({ setIsModalOpen, data, id, currentType }: EditContactsProps) => {
  const dispatch = useAppDispatch();
  const [isProcessing, setIsProcessing] = useState(false);
  const navigate = useNavigate();
  const {
    handleSubmit,
    control,
    setValue,
    formState: { errors, isDirty, isValid }
  } = useForm<ContactsFormInput>({
    resolver: zodResolver(contactsValidation),
    mode: 'onChange',
    defaultValues: {}
  });

  useEffect(() => {
    setValue(currentType as keyof ContactsFormInput, data);
  }, [currentType, data, setValue]);

  const onSubmit: SubmitHandler<ContactsFormInput> = async (
    values: ContactsFormInput
  ) => {
    try {
      setIsProcessing(true);
      if (currentType === 'email') {
        await dispatch(
          editEmail({
            id,
            values: { email: values.email || '' }
          })
        );
      } else if (currentType === 'phone') {
        await dispatch(
          editPhone({
            id,
            values: { phone: values.phone || '' }
          })
        );
      } else if (currentType === 'name') {
        await dispatch(
          editName({
            id,
            values: { name: values.name || '' }
          })
        );
      }
      setIsProcessing(false);
      dispatch(openAlert(editSuccessResponseMessage('контакт')));
      setIsModalOpen(false);
      navigate(0);
    } catch (error: unknown) {
      console.error('Error updating contact:', error);
      dispatch(openAlert(editErrorResponseMessage('контакт')));
    }
  };

  console.log(data);
  console.log(currentType);

  return (
    <div className="left-1/6 fixed top-0 z-20 h-full w-5/6 bg-[rgba(0,0,0,0.6)]">
      <div className="absolute left-[50%] top-[50%] z-[9999] flex h-[60vh] w-[50vw] -translate-x-[50%] -translate-y-[50%] items-center justify-center gap-4 bg-white px-4 py-8 text-black">
        <button
          className="absolute right-5 top-4 text-graphite hover:text-accent"
          onClick={() => setIsModalOpen(false)}
        />
        <form
          onSubmit={handleSubmit(onSubmit)}
          autoComplete="off"
          className="mx-auto w-full max-w-md"
        >
          <div className="flex flex-col gap-6 p-6">
            <div className="flex items-center justify-between">
              <h4 className="text-2xl font-bold">
                {`Зміна ${
                  currentType === 'email'
                    ? 'електронної пошти'
                    : currentType === 'phone'
                    ? 'номера телефону'
                    : "ім'я"
                }`}
              </h4>
              <button
                type="button"
                onClick={() => setIsModalOpen(false)}
                className="text-gray-500 hover:text-gray-700"
              >
                <TfiClose size={24} />
              </button>
            </div>

            <p className="text-graphite">
              {`Ваш ${
                currentType === 'email'
                  ? 'електронна пошта'
                  : currentType === 'phone'
                  ? 'номер телефону'
                  : "ім'я"
              }: `}
              <span className="text-[17px] font-medium">{data}</span>
            </p>

            <div className="mt-4">
              <Controller
                control={control}
                name={
                  currentType === 'email'
                    ? 'email'
                    : currentType === 'phone'
                    ? 'phone'
                    : 'name'
                }
                render={({ field }) => (
                  <TextInput
                    {...field}
                    errorText={
                      currentType === 'email'
                        ? errors.email?.message
                        : currentType === 'phone'
                        ? errors.phone?.message
                        : errors.name?.message
                    }
                    placeholder={
                      currentType === 'email'
                        ? 'Введіть електронну пошту'
                        : currentType === 'phone'
                        ? 'Введіть номер телефону'
                        : "Введіть ім'я"
                    }
                    title={
                      currentType === 'email'
                        ? 'Електронна пошта'
                        : currentType === 'phone'
                        ? 'Номер телефону'
                        : "Ім'я"
                    }
                  />
                )}
              />
            </div>

            <p
              className={`text-[17px] ${
                isDirty && isValid ? 'text-black' : 'text-gray-400'
              }`}
            >
              {`Змінити ${
                currentType === 'email'
                  ? 'електронну пошту'
                  : currentType === 'phone'
                  ? 'номер телефону'
                  : "ім'я"
              }?`}
            </p>

            <div className="mt-4 flex gap-4">
              <button
                type="submit"
                disabled={!isDirty || !isValid || isProcessing}
                className={`flex-1 rounded px-6 py-2 font-medium ${
                  isDirty && isValid && !isProcessing
                    ? 'bg-accent text-white hover:bg-accent/90'
                    : 'cursor-not-allowed bg-gray-200 text-gray-500'
                }`}
              >
                {isProcessing ? 'Обробка запиту...' : 'Зберегти зміни'}
              </button>
              <button
                type="button"
                onClick={() => setIsModalOpen(false)}
                className="flex-1 rounded border border-gray-300 bg-white px-6 py-2 text-gray-700 hover:bg-gray-50"
              >
                Скасувати
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Edit;
