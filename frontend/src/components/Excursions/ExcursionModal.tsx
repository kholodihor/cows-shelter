import people_icon from '@/assets/icons/people_icon.svg';
import time_icon from '@/assets/icons/time_icon.svg';
import close_icon_black from '@/assets/icons/close_icon_black.svg';
import { closeModal, openModal } from '@/store/slices/modalSlice';
import {
  Dispatch,
  MouseEvent,
  SetStateAction,
  useEffect,
  useState
} from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hook';
import { useTranslation } from 'react-i18next';

type ExcursionModalProps = {
  isOpen: boolean;
  gallery: Array<{
    id: string;
    image_url: string;
    [key: string]: any;
  }>;
  setShowModal: Dispatch<SetStateAction<boolean>>;
};

type ExcursionData = {
  id: string;
  title_ua: string;
  title_en: string;
  description_ua: string;
  description_en: string;
  amount_of_persons: string;
  time_from: string;
  time_to: string;
  image_url: string;
  [key: string]: any; // For any additional properties
};

const ExcursionModal = ({
  isOpen,
  gallery,
  setShowModal
}: ExcursionModalProps) => {
  const { i18n } = useTranslation();
  const { language } = i18n;
  const dispatch = useAppDispatch();
  const excursion = useAppSelector(
    (state) => state.modals.data
  ) as ExcursionData;
  const [imagesToDisplay, setImagesToDisplay] = useState<
    Array<{ id: string; image_url: string }>
  >([]);

  const [windowWidth, setWindowWidth] = useState(window.innerWidth);
  const [windowHeight, setHeight] = useState(window.innerHeight);

  useEffect(() => {
    const generateUniqueRandomIndices = (
      galleryLength: number,
      count: number
    ) => {
      const indices = new Set<number>();

      while (indices.size < count) {
        const randomIndex = generateRandom(galleryLength - 1);
        indices.add(randomIndex); // Add to the set, which ensures uniqueness
      }

      return Array.from(indices); // Convert the set to an array
    };

    if (gallery.length >= 3) {
      const uniqueRandomIndices = generateUniqueRandomIndices(
        gallery.length,
        3
      );
      const randomImages = uniqueRandomIndices.map((index) => gallery[index]);

      setImagesToDisplay(randomImages);
    }
  }, [gallery]);

  const handleChangedSize = () => {
    setWindowWidth(window.innerWidth);
    setHeight(window.innerHeight);
  };

  useEffect(() => {
    window.addEventListener('resize', handleChangedSize);
    return () => {
      window.removeEventListener('resize', handleChangedSize);
    };
  }, []);

  const closeExcursionsModal = () => {
    setShowModal(false);
    setTimeout(() => {
      dispatch(closeModal());
    }, 500);
  };

  const handleOverlayClick = (event: MouseEvent<HTMLDivElement>) => {
    if (event.currentTarget === event.target) {
      closeExcursionsModal();
    }
  };

  const openExcursionOrderModal = () => {
    setTimeout(() => {
      dispatch(openModal({ data: {}, type: 'order' }));
    }, 300);
  };

  const openDonationModal = () => {
    dispatch(openModal({ data: {}, type: 'donation' }));
  };

  function generateRandom(max: number) {
    let rand = Math.random() * max;
    rand = Math.floor(rand);
    return rand;
  }

  return (
    <div
      className={`fixed left-0 top-0 z-50 h-screen w-screen bg-black transition-all duration-700 ${
        isOpen ? 'bg-opacity-40' : 'bg-opacity-0'
      } `}
      onClick={handleOverlayClick}
    >
      <div
        className={`absolute left-1/2 top-1/2 ${
          isOpen ? 'translate-x-0' : 'translate-x-[100%]'
        } w-[85%] ${
          windowHeight < 400 && 'md:h-[380px]'
        } h-[675px] max-h-full translate-x-[-50%] translate-y-[-50%] overflow-auto bg-white 
        px-5 pb-12 pt-4 transition-all duration-700 md:h-[832px] md:w-[672px] md:p-10 lg:h-auto lg:w-[1136px] 
         lg:px-[3.75rem] lg:pb-[3.75rem]`}
      >
        <div className="flex flex-col md:gap-6 lg:flex-row lg:gap-10">
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-2 md:flex-row md:gap-6 lg:gap-9">
              <div className="flex gap-3">
                <img src={time_icon} alt="Time icon" width={24} height={24} />
                <span className="text-sm leading-normal text-darkyellow md:text-base">
                  {excursion?.time_from} - {excursion?.time_to}{' '}
                  {language === 'uk' ? 'хвилин' : 'minutes'}
                </span>
              </div>
              <div className="flex gap-3">
                <img
                  src={people_icon}
                  alt="People icon"
                  width={24}
                  height={24}
                />
                <span className="text-sm leading-normal text-darkyellow md:text-base">
                  {language === 'uk' ? 'до' : 'up to'}{' '}
                  {excursion?.amount_of_persons}{' '}
                  {language === 'uk' ? 'відвідувачів' : 'visitors'}
                </span>
              </div>
            </div>
            {windowWidth >= 1280 && (
              <div className="grid grid-cols-2 gap-3.5">
                {imagesToDisplay.length > 0 ? (
                  imagesToDisplay.map((item, index) => (
                    <div
                      key={`${item.id}-${index}`}
                      className={`group relative cursor-pointer ${
                        index === 2
                          ? 'col-span-2 h-[260px] w-[490px]'
                          : 'h-[260px] w-[238px]'
                      }`}
                    >
                      <img
                        src={item.image_url}
                        alt={`Excursion ${index + 1}`}
                        className="h-full w-full object-cover"
                        onError={(e) => {
                          const target = e.target as HTMLImageElement;
                          target.src = '/placeholder-image.jpg'; // Fallback image
                        }}
                      />
                    </div>
                  ))
                ) : (
                  <div className="col-span-2 flex h-64 items-center justify-center bg-gray-100">
                    <p className="text-gray-500">No images available</p>
                  </div>
                )}
              </div>
            )}
            {windowWidth >= 768 && windowWidth < 1280 && (
              <div className="flex gap-3">
                <img
                  src={imagesToDisplay[0]?.image_url}
                  className="h-[194px] w-[290px] object-cover"
                ></img>
                <img
                  src={imagesToDisplay[1]?.image_url}
                  className="h-[194px] w-[290px] object-cover"
                ></img>
              </div>
            )}
            {windowWidth < 768 && (
              <img
                src={imagesToDisplay[0]?.image_url}
                className="min-h-[186px] min-w-[280px] object-cover"
              ></img>
            )}
          </div>
          <div className="mt-[1.25rem] flex w-[90%] flex-col justify-between gap-3.5 md:mt-0 md:w-[592px] md:gap-6 lg:w-[28.75rem]">
            <h3 className="pb-[7px] text-lg font-semibold leading-normal text-black md:mt-0 md:text-xl lg:border-b lg:border-disabled lg:text-2xl lg:font-bold">
              {language === 'uk' ? excursion?.title_ua : excursion?.title_en}
            </h3>
            <p className="whitespace-pre-line text-sm font-normal leading-normal text-black md:text-base">
              {language === 'uk'
                ? excursion?.description_ua
                : excursion?.description_en}
            </p>
            <div className="flex flex-col gap-3 md:flex-row md:justify-center md:gap-6">
              <button
                className="h-10 max-w-[17.5rem] bg-accent text-lg font-medium leading-[1.375rem] transition-all duration-300 hover:bg-lemon focus:bg-lemon active:bg-lemon md:h-11 md:w-[14.44rem]"
                onClick={openExcursionOrderModal}
              >
                {language === 'uk'
                  ? 'Замовити екскурсію'
                  : 'Order an excursion'}
              </button>
              <button
                className="h-10 max-w-[17.5rem] border border-solid border-black text-lg font-medium leading-[1.375rem] text-black hover:border-accent focus:border-accent active:border-accent md:h-11 md:w-[14.44rem]"
                onClick={openDonationModal}
              >
                {language === 'uk' ? 'Допомогти' : 'Donate'}
              </button>
            </div>
          </div>
        </div>
        <button
          className="absolute right-6 top-6"
          onClick={() => closeExcursionsModal()}
        >
          <img src={close_icon_black} width={44} height={44}></img>
        </button>
      </div>
    </div>
  );
};

export default ExcursionModal;
