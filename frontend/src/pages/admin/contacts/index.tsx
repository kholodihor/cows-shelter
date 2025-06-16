import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import Pen from '@/components/icons/Pen';
import PlusIcon from '@/components/icons/PlusIcon';
import { useAppDispatch, useAppSelector } from '@/store/hook';
import Edit from './edit';
import { fetchContacts } from '@/store/slices/contactsSlice';
import Loader from '@/components/admin/Loader';
import ResponseAlert from '@/components/admin/ResponseAlert';

interface Contact {
  id: number;
  name: string;
  email: string;
  phone: string;
}

const Contacts = () => {
  const dispatch = useAppDispatch();
  const [modalData, setModalData] = useState('');
  const [currentEditType, setCurrentEditType] = useState<
    'name' | 'email' | 'phone'
  >('email');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const contacts = useAppSelector(
    (state) => state.contacts.contacts
  ) as Contact[];
  const isLoading = useAppSelector((state) => state.contacts.loading);
  const error = useAppSelector((state) => state.contacts.error);
  const isAlertOpen =
    useAppSelector((state) => state.alert?.isAlertOpen) || false;

  useEffect(() => {
    dispatch(fetchContacts());
  }, [dispatch, isModalOpen]);

  // Show loader while loading
  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Loader />
      </div>
    );
  }

  // Show error message if there's an error
  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-red-500">Помилка завантаження контактів: {error}</p>
      </div>
    );
  }

  // Check if contacts exist and has at least one item
  const contact = contacts?.[0];

  // If no contacts found, show the add button
  if (!contact) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center">
        <p className="mb-4">Контакти не знайдено</p>
        <Link
          to="/admin/contacts/add"
          className="flex items-center gap-2 rounded bg-accent px-4 py-2 text-white transition-colors hover:bg-accent/90"
        >
          <PlusIcon className="h-5 w-5" />
          Додати контакти
        </Link>
      </div>
    );
  }

  // If contact exists, show the contact details without the add button
  return (
    <div className="relative flex min-h-screen flex-col items-start pt-10">
      <div className="mb-8 px-8">
        <h3 className="text-[32px] font-semibold">Контакти</h3>
      </div>

      <div className="mb-4 ml-8 flex max-w-[21.5rem] flex-col items-center justify-between border border-black px-[90px] pb-2 pt-8">
        <h2 className="pb-2 text-center text-base">Ваше ім'я:</h2>
        <p className="pb-8 text-lg font-bold text-black">
          {contact?.name || 'Не вказано'}
        </p>
        <div className="flex w-full">
          <button
            onClick={() => {
              setModalData(contact?.name || '');
              setCurrentEditType('name');
              setIsModalOpen(true);
            }}
            className="w-full border border-darkgray bg-lightgrey px-[4rem] py-2 text-center text-xl text-black hover:text-accent"
          >
            <Pen />
          </button>
        </div>
      </div>

      <div className="mb-4 ml-8 flex max-w-[21.5rem] flex-col items-center justify-between border border-black px-[90px] pb-2 pt-8">
        <h2 className="pb-2 text-center text-base">Ваша електронна адреса:</h2>
        <p className="pb-8 text-lg font-bold text-black">
          {contact?.email || 'Не вказано'}
        </p>
        <div className="flex w-full">
          <button
            onClick={() => {
              setModalData(contact?.email || '');
              setCurrentEditType('email');
              setIsModalOpen(true);
            }}
            className="w-full border border-darkgray bg-lightgrey px-[4rem] py-2 text-center text-xl text-black hover:text-accent"
          >
            <Pen />
          </button>
        </div>
      </div>

      <div className="ml-8 flex max-w-[21.5rem] flex-col items-center justify-between border border-black px-[90px] pb-2 pt-8">
        <h2 className="pb-2 text-center text-base">Ваш номер телефону:</h2>
        <p className="pb-8 text-lg font-bold text-black">
          {contact?.phone || 'Не вказано'}
        </p>
        <div className="flex w-full bg-lightgrey">
          <button
            onClick={() => {
              setModalData(contact?.phone || '');
              setCurrentEditType('phone');
              setIsModalOpen(true);
            }}
            className="w-full border border-darkgray px-[4.65rem] py-2 text-xl text-black hover:text-accent"
          >
            <Pen />
          </button>
        </div>
      </div>

      {isModalOpen && (
        <Edit
          setIsModalOpen={setIsModalOpen}
          data={modalData}
          id={contact?.id.toString()}
          currentType={currentEditType}
        />
      )}

      {isAlertOpen && <ResponseAlert />}
    </div>
  );
};

export default Contacts;
