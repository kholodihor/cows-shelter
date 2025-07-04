import { useEffect, useState } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hook';
import { openModal } from '@/store/slices/modalSlice';
import { useTranslation } from 'react-i18next';
import { useInView } from 'react-intersection-observer';
import { images } from '@/data/gallery';

import ExcursionModal from './ExcursionModal';
import LittleArrow from '../icons/LittleArrow';
import Slider from '@/components/Slider';
import { setActiveLink } from '@/store/slices/observationSlice';
import { Excursion, fetchExcursion } from '@/store/slices/excursionsSlice';
import { useWidth } from '@/hooks/useWidth';
import LoadingSpinner from '../common/LoadingSpinner';
import ErrorMessage from '../common/ErrorMessage';

const Excursions = () => {
  const { i18n } = useTranslation();
  const { language } = i18n;
  const dispatch = useAppDispatch();
  const screenWidth = useWidth();
  const [showModal, setShowModal] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(6);
  const [pagesLength, setPagesLength] = useState(0);
  const type = useAppSelector((state) => state.modals.type);
  const isModalOpen = useAppSelector((state) => state.modals.isModalOpen);

  // Get excursions data from Redux store
  const { excursions, loading, error } = useAppSelector((state) => ({
    excursions: state.excursions.excursions,
    loading: state.excursions.loading,
    error: state.excursions.error
  }));

  const totalLength = excursions.length;

  console.log(currentPage);

  const { ref, inView } = useInView({
    threshold: 0.5
  });

  // Fetch excursions data on component mount
  useEffect(() => {
    dispatch(fetchExcursion());
  }, [dispatch]);

  useEffect(() => {
    if (screenWidth >= 1280) {
      setItemsPerPage(3);
    } else if (screenWidth >= 768) {
      setItemsPerPage(2);
    } else {
      setItemsPerPage(totalLength);
    }
  }, [screenWidth, totalLength]);

  useEffect(() => {
    const pagesNumber = totalLength / itemsPerPage;
    setPagesLength(pagesNumber < 5 ? pagesNumber : 5);
  }, [totalLength, itemsPerPage]);

  const openExcursionModal = (item: Excursion) => {
    dispatch(openModal({ data: item, type: 'excursions' }));
  };

  useEffect(() => {
    if (isModalOpen) {
      setShowModal(true);
    } else {
      setShowModal(false);
    }
  }, [isModalOpen]);

  useEffect(() => {
    if (inView) {
      dispatch(setActiveLink('#excursions'));
    } else {
      dispatch(setActiveLink(''));
    }
  }, [inView, dispatch]);

  const [windowWidth, setWindowWidth] = useState(window.innerWidth);

  const handleChangedSize = () => {
    setWindowWidth(window.innerWidth);
  };

  useEffect(() => {
    window.addEventListener('resize', handleChangedSize);
    return () => {
      window.removeEventListener('resize', handleChangedSize);
    };
  }, []);

  // Show loading state
  if (loading) {
    return (
      <section
        id="excursions"
        className="relative bg-[#F3F3F5] py-6 md:py-16 lg:py-20"
      >
        <div className="mx-auto px-5 sm:w-[480px] md:w-[768px] md:px-12 lg:w-[1280px] lg:px-[120px] xl:w-[1440px]">
          <div className="flex h-64 items-center justify-center">
            <LoadingSpinner />
          </div>
        </div>
      </section>
    );
  }

  // Show error state
  if (error) {
    return (
      <section
        id="excursions"
        className="relative bg-[#F3F3F5] py-6 md:py-16 lg:py-20"
      >
        <div className="mx-auto px-5 sm:w-[480px] md:w-[768px] md:px-12 lg:w-[1280px] lg:px-[120px] xl:w-[1440px]">
          <ErrorMessage
            message={error}
            onRetry={() => dispatch(fetchExcursion())}
          />
        </div>
      </section>
    );
  }

  // Show empty state if no excursions
  if (excursions.length === 0) {
    return (
      <section
        id="excursions"
        className="relative bg-[#F3F3F5] py-6 md:py-16 lg:py-20"
      >
        <div className="mx-auto px-5 sm:w-[480px] md:w-[768px] md:px-12 lg:w-[1280px] lg:px-[120px] xl:w-[1440px]">
          <p className="text-center text-lg">
            {language === 'uk'
              ? 'Наразі екскурсії відсутні'
              : 'No excursions available at the moment'}
          </p>
        </div>
      </section>
    );
  }

  return (
    <section
      id="excursions"
      ref={ref}
      className="relative bg-[#F3F3F5] py-6 md:py-16 lg:py-20"
    >
      <div className="mx-auto px-5 sm:w-[480px] md:w-[768px] md:px-12 lg:w-[1280px] lg:px-[120px] xl:w-[1440px]">
        {windowWidth < 768 && (
          <div>
            <h2 className="mb-10 flex gap-2 text-2xl font-medium leading-normal">
              {language === 'uk' ? 'Екскурсії' : 'Excursions'}
              <img src="/cow.svg" alt="" className="" />
            </h2>
            <ul className="mb-5 flex flex-col items-center gap-4">
              {excursions && Array.isArray(excursions) ? (
                excursions.map((excursion: (typeof excursions)[0]) => (
                  <li key={excursion.ID} className="drop-shadow">
                    <div className="relative">
                      <img
                        src={
                          excursion.image_url || '/placeholder-excursion.jpg'
                        }
                        alt={
                          language === 'uk'
                            ? excursion.title_ua
                            : excursion.title_en
                        }
                        width={280}
                        height={200}
                        className="h-full w-full object-cover"
                      />
                      <div className="absolute bottom-0 left-0 flex flex-col gap-3 pb-5 pl-5 text-white">
                        <p className="text-lg font-medium leading-normal">
                          {language === 'uk'
                            ? excursion.title_ua
                            : excursion.title_en}
                        </p>
                        <a>
                          <button
                            className="flex gap-3 border border-solid border-white py-2.5 pl-5 pr-4"
                            onClick={() => {
                              openExcursionModal(excursion as any);
                            }}
                          >
                            <span className="text-base leading-tight">
                              {language === 'uk'
                                ? 'Показати більше'
                                : 'Show more'}
                            </span>
                            <LittleArrow />
                          </button>
                        </a>
                      </div>
                    </div>
                  </li>
                ))
              ) : (
                <p className="text-black">Сервер не відповідає</p>
              )}
            </ul>
          </div>
        )}
        {windowWidth >= 768 && (
          <Slider
            title={language === 'uk' ? 'Екскурсії' : 'Excursions'}
            setCurrentPage={setCurrentPage}
            pagesLength={pagesLength}
            isExcursions={true}
          >
            <ul className="mb-0 flex justify-center gap-6 lg:mb-[8.75rem]">
              {excursions && Array.isArray(excursions) ? (
                excursions.map((excursion: (typeof excursions)[0]) => (
                  <li key={excursion.ID} className="drop-shadow">
                    <div className="group relative">
                      <img
                        className="h-[235px] w-[324px] object-cover lg:h-[384px] lg:w-[384px]"
                        src={
                          excursion.image_url || '/placeholder-excursion.jpg'
                        }
                        alt={
                          language === 'uk'
                            ? excursion.title_ua
                            : excursion.title_en
                        }
                      />
                      <div className="absolute left-0 top-0 h-full w-full lg:bg-black/[.60] lg:opacity-0 lg:transition-all lg:duration-700 lg:group-hover:opacity-100"></div>
                      <div className="absolute bottom-0 left-0 flex flex-col gap-0 pb-6 pl-6 text-white transition-all duration-700 group-hover:gap-5">
                        <p className="text-xl font-semibold leading-normal lg:text-[22px]">
                          {language === 'uk'
                            ? excursion.title_ua
                            : excursion.title_en}
                        </p>
                        <p className="opacity-0 transition-all duration-700 lg:group-hover:opacity-100">
                          {excursion.time_from} - {excursion.time_to}{' '}
                          {language === 'uk' ? 'хвилин' : 'minutes'} /{' '}
                          {language === 'uk' ? 'до' : 'up to'}{' '}
                          {excursion.amount_of_persons}{' '}
                          {language === 'uk' ? 'відвідувачів' : 'visitors'}
                        </p>
                        <a>
                          <button
                            className="flex gap-3 border border-solid border-white py-[0.69rem] pl-6 pr-2.5 transition-all duration-700 lg:border-transparent lg:group-hover:border-white lg:group-hover:hover:border-accent lg:focus:border-accent lg:active:border-accent"
                            onClick={() => {
                              openExcursionModal(excursion as any);
                            }}
                          >
                            <span className="text-lg font-medium leading-tight">
                              {language === 'uk'
                                ? 'Показати більше'
                                : 'Show more'}
                            </span>
                            <LittleArrow />
                          </button>
                        </a>
                      </div>
                    </div>
                  </li>
                ))
              ) : (
                <p className="text-black">Сервер не відповідає</p>
              )}
            </ul>
          </Slider>
        )}
        {isModalOpen && type === 'excursions' && (
          <ExcursionModal
            isOpen={showModal}
            setShowModal={setShowModal}
            gallery={images}
          />
        )}
      </div>
    </section>
  );
};

export default Excursions;
