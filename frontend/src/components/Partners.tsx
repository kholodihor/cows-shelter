import { useTranslation } from 'react-i18next';
import { useEffect, useState } from 'react';
import { useWidth } from '@/hooks/useWidth';
import { useInView } from 'react-intersection-observer';
import { setActiveLink } from '@/store/slices/observationSlice';
import { useAppDispatch, useAppSelector } from '@/store/hook';

// Redux actions
import { fetchPartners } from '@/store/slices/partnersSlice';

// Components
import Slider from '@/components/Slider';
import PartnersModal from './modals/PartnersModal';
import { openModal } from '@/store/slices/modalSlice';
import Loader from './admin/Loader';
import ErrorMessage from './common/ErrorMessage';

const Partners = () => {
  const { t } = useTranslation();
  const screenWidth = useWidth();
  const dispatch = useAppDispatch();

  const [pagesLength, setPagesLength] = useState(0);
  const [itemsPerPage, setItemsPerPage] = useState(4);
  // currentPage is not currently used but kept for future pagination
  const [_currentPage, setCurrentPage] = useState(1);
  const [showModal, setShowModal] = useState(false);

  // Get partners data from Redux store
  const { partners, loading, error } = useAppSelector((state) => ({
    partners: state.partners.partners,
    loading: state.partners.loading,
    error: state.partners.error
  }));

  const type = useAppSelector((state) => state.modals.type);
  const isModalOpen = useAppSelector((state) => state.modals.isModalOpen);
  const totalLength = partners.length;

  const { ref, inView } = useInView({
    threshold: 0.5
  });

  // Fetch partners data on component mount
  useEffect(() => {
    dispatch(fetchPartners())
      .unwrap()
      .catch((error) => {
        console.error('Failed to fetch partners:', error);
      });
  }, [dispatch]);

  useEffect(() => {
    if (screenWidth > 1280) {
      setItemsPerPage(4);
    }
    if (screenWidth >= 768 && screenWidth < 1280) {
      setItemsPerPage(3);
    }
    if (screenWidth > 320 && screenWidth < 768) {
      setItemsPerPage(totalLength);
    }
  }, [screenWidth, totalLength]);

  useEffect(() => {
    const pagesNumber = totalLength / itemsPerPage;
    setPagesLength(pagesNumber < 5 ? pagesNumber : 5);
  }, [totalLength, itemsPerPage]);

  const openPartnersModal = () => {
    dispatch(openModal({ data: {}, type: 'partners' }));
  };

  useEffect(() => {
    if (isModalOpen) {
      setTimeout(() => {
        setShowModal(true);
      }, 300);
    } else {
      setTimeout(() => {
        setShowModal(false);
      }, 300);
    }
  }, [isModalOpen]);

  useEffect(() => {
    if (inView) {
      dispatch(setActiveLink('#partners'));
    } else {
      dispatch(setActiveLink(''));
    }
  }, [inView, dispatch]);

  // Show loading state
  if (loading) {
    return (
      <section id="partners" className="relative bg-[#F3F3F5] py-12">
        <div className="mx-auto px-5 sm:w-[480px] md:w-[768px] md:px-8 lg:w-[1280px] lg:px-[120px] lg:py-[80px]">
          <div className="flex h-64 items-center justify-center">
            <Loader />
          </div>
        </div>
      </section>
    );
  }

  // Show error state
  if (error) {
    return (
      <section id="partners" className="relative bg-[#F3F3F5] py-12">
        <div className="mx-auto px-5 sm:w-[480px] md:w-[768px] md:px-8 lg:w-[1280px] lg:px-[120px] lg:py-[80px]">
          <ErrorMessage
            message={error}
            onRetry={() => dispatch(fetchPartners())}
          />
        </div>
      </section>
    );
  }

  // Show empty state if no partners
  if (partners.length === 0) {
    return (
      <section id="partners" className="relative bg-[#F3F3F5] py-12">
        <div className="mx-auto px-5 sm:w-[480px] md:w-[768px] md:px-8 lg:w-[1280px] lg:px-[120px] lg:py-[80px]">
          <div className="sectionHeader mb-5 flex-row md:mb-8 lg:mb-14 lg:flex lg:items-center lg:justify-between">
            <h2 className="flex gap-2 text-[1.5rem] font-medium md:mb-6 md:text-[44px] lg:text-[54px]">
              {t('partners:header')}
              <img src="/cow.svg" alt="" className="" />
            </h2>
          </div>
          <p className="text-center text-lg">
            {t('partners:noPartners', 'No partners available')}
          </p>
        </div>
      </section>
    );
  }

  return (
    <section id="partners" ref={ref} className="relative bg-[#F3F3F5] ">
      <div className="mx-auto flex flex-col px-5 py-6 sm:w-[480px] md:w-[768px] md:p-12 lg:w-[1280px] lg:px-[120px] lg:py-[80px] xl:w-[1440px]">
        <div className="sectionHeader mb-5 flex-row md:mb-8 lg:mb-14 lg:flex  lg:items-center lg:justify-between">
          <h2 className="flex gap-2 text-[1.5rem] font-medium md:mb-6 md:text-[44px] lg:text-[54px]">
            {t('partners:header')}
            <img src="/cow.svg" alt="" className="" />
          </h2>
          <button
            onClick={openPartnersModal}
            className="duration-800 hidden bg-accent px-8 py-3 text-lg font-medium leading-6 hover:bg-lemon focus:bg-lemon active:bg-darkyellow md:inline-block"
          >
            {t('partners:become_partner')}
          </button>
        </div>
        <p className="mb-5 text-[18px] leading-relaxed text-gray-700  md:mb-10 md:text-[20px] lg:w-[1070px] lg:text-[22px]">
          {t('partners:text')}
        </p>
        {/* mobile only */}
        <ul className="mb-5 grid grid-cols-2 gap-x-3 gap-y-2.5 overflow-x-auto md:hidden ">
          {partners && Array.isArray(partners) ? (
            partners.map((partner) => (
              <li
                key={partner?.ID}
                className="flex flex-col   md:mb-6 md:w-[calc(50%-0.75rem)]"
              >
                <a
                  href={partner.link}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="partner-scale block w-full border-solid border-darkyellow transition-all duration-300 hover:border-b"
                >
                  <img
                    className="m-auto mb-6 scale-100"
                    src={partner.logo || '/placeholder-partner.png'}
                    alt={partner.name}
                    width={134}
                    height={134}
                    loading="lazy"
                  />
                  <p className="mb-4 text-center text-[1rem] leading-relaxed md:text-[20px] lg:text-[22px]">
                    {partner.name}
                  </p>
                </a>
              </li>
            ))
          ) : (
            <p>Сервер не відповідає</p>
          )}
        </ul>
        {/* tablet + desktop */}
        <div className="hidden md:block">
          <Slider
            isPartners
            setCurrentPage={setCurrentPage}
            pagesLength={pagesLength}
          >
            <ul className="mb-5  flex gap-6 lg:mt-20">
              {partners && Array.isArray(partners) ? (
                partners.map((partner) => (
                  <li key={partner.ID} className="flex justify-around">
                    <a
                      href={partner.link}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="partner-scale block w-full border-solid border-darkyellow transition-all duration-300 hover:border-b"
                    >
                      <img
                        className="m-auto mb-6 scale-100 object-contain md:h-[208px] md:w-[208px] xl:h-[245px] xl:w-[245px]"
                        src={partner.logo || '/placeholder-partner.png'}
                        alt={partner.name}
                        width={208}
                        height={208}
                        loading="lazy"
                      />
                      <p className=":w-[282px] mb-5 w-[208px] text-center text-[1rem] leading-relaxed md:text-[20px] lg:mb-[4.125rem] lg:text-[22px]">
                        {partner.name}
                      </p>
                    </a>
                  </li>
                ))
              ) : (
                <p>Сервер не відповідає</p>
              )}
            </ul>
          </Slider>
        </div>{' '}
        <button
          onClick={openPartnersModal}
          className="duration-800 bg-accent  px-8 py-3 text-lg font-medium leading-6 hover:bg-lemon focus:bg-lemon active:bg-darkyellow md:hidden"
        >
          {t('partners:become_partner')}
        </button>
      </div>{' '}
      {isModalOpen && type === 'partners' && (
        <PartnersModal isOpen={showModal} setShowModal={setShowModal} />
      )}
    </section>
  );
};

export default Partners;
