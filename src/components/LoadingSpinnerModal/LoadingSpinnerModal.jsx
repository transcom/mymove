import React from 'react';
import PropTypes from 'prop-types';
import { Oval } from 'react-loader-spinner';

import styles from './LoadingSpinnerModal.module.scss';

import Modal, { ModalTitle, connectModal } from 'components/Modal/Modal';

const LoadingSpinnerModal = ({ message }) => {
  return (
    <Modal>
      <ModalTitle className={styles.center}>
        <h3>{message}</h3>
      </ModalTitle>
      <div className={styles.center}>
        <Oval
          visible
          height="100"
          width="100"
          color="#ffbe2e"
          secondaryColor="#252f3e"
          ariaLabel="oval-loading"
          wrapperStyle={{}}
          wrapperClass=""
        />
      </div>
    </Modal>
  );
};

LoadingSpinnerModal.propTypes = {
  message: PropTypes.string,
};

LoadingSpinnerModal.defaultProps = {
  message: '',
};

export default connectModal(LoadingSpinnerModal);
