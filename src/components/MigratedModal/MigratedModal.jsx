/* eslint-disable react/jsx-props-no-spreading */
import React, { useEffect, useState, useRef } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './MigratedModal.module.scss';

const trapFocus = (modalRef) => {
  const focusableSelectors = [
    'a[href]',
    'button:not([disabled])',
    'textarea:not([disabled])',
    'input:not([disabled])',
    '[tabindex]:not([tabindex="-1"])',
  ];
  const focusableElements = modalRef.current?.querySelectorAll(focusableSelectors.join(','));
  return Array.prototype.slice.call(focusableElements);
};

/** This is a straightforward port of the Modal component from React-USWDS 1.17
 *  into the MilMove project, as the component is being deprecated in USWDS 2.x. */

/** Modal UI component */
export const Modal = ({ title, children, actions, className, onClose }) => {
  const classes = classnames(styles.modal, className);
  const modalRef = useRef(null);

  useEffect(() => {
    if (modalRef.current) {
      modalRef.current.focus();
    }
  }, []);

  useEffect(() => {
    const handleKeyDown = (event) => {
      /** Ignores hidden modals in the DOM */
      const modal = modalRef.current;
      if (!modal || modal.offsetParent === null) return;

      if (event.key === 'Escape') {
        onClose?.();
        return;
      }

      if (event.key === 'Tab' && modalRef.current) {
        const focusableEls = trapFocus(modalRef);
        const firstEl = focusableEls[0];
        const lastEl = focusableEls[focusableEls.length - 1];
        const activeIndex = focusableEls.indexOf(document.activeElement);

        if (activeIndex === -1) {
          event.preventDefault();
          lastEl?.focus();
          return;
        }

        if (!event.shiftKey && document.activeElement === lastEl) {
          event.preventDefault();
          firstEl?.focus();
        }

        if (event.shiftKey && document.activeElement === firstEl) {
          event.preventDefault();
          lastEl?.focus();
        }
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [onClose]);

  return (
    <div data-testid="modal" className={classes} ref={modalRef}>
      <div className={styles.title}>{title}</div>
      <div className={styles.content}>{children}</div>
      <div className={styles.actions}>{actions}</div>
    </div>
  );
};

Modal.propTypes = {
  title: PropTypes.node,
  children: PropTypes.node.isRequired,
  actions: PropTypes.node,
  className: PropTypes.string,
  onClose: PropTypes.func,
};

Modal.defaultProps = {
  title: '',
  actions: '',
  className: '',
  onClose: () => {},
};

/** Overlay UI component (grey background) */
export const Overlay = () => <div className={styles.overlay} />;

/** Modal positioning component */
export const ModalContainer = ({ children }) => {
  return <div className={styles.modalContainer}>{children}</div>;
};

ModalContainer.propTypes = {
  children: PropTypes.node.isRequired,
};

export const connectModal = (Component) => {
  const ConnectedModal = ({ isOpen, ...props }) => {
    if (!isOpen) return null;
    return (
      <>
        <Overlay />
        <ModalContainer>
          <Component {...props} />
        </ModalContainer>
      </>
    );
  };

  ConnectedModal.propTypes = {
    isOpen: PropTypes.bool,
  };

  ConnectedModal.defaultProps = {
    isOpen: false,
  };

  return ConnectedModal;
};

export const useModal = () => {
  const [isOpen, setIsOpen] = useState(false);

  const openModal = () => {
    setIsOpen(true);
  };
  const closeModal = () => {
    setIsOpen(false);
  };

  return { isOpen, openModal, closeModal };
};
