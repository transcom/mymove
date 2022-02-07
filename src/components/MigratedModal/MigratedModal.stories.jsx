import React from 'react';
import { Button } from '@trussworks/react-uswds';

import { Modal, Overlay, ModalContainer, connectModal, useModal } from './MigratedModal';

export default {
  title: 'Components/MigratedModal',
};

/** This is a straightforward port of the Modal component from React-USWDS 1.17
 *  into the MilMove project, as the component is being deprecated in USWDS 2.x. */

const testTitle = <h2>My Modal</h2>;

const testActions = (
  <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
    <Button type="button" secondary>
      Cancel
    </Button>
    <Button type="button">Close</Button>
  </div>
);

const closeCaseActions = (
  <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
    <Button type="button" secondary>
      Cancel
    </Button>
    <Button type="button">Close Investigation</Button>
  </div>
);

export const basic = () => (
  <div
    style={{
      margin: '100px',
    }}
  >
    <Modal>
      <p>This is a basic modal!</p>
    </Modal>
  </div>
);

export const withTitle = () => (
  <div
    style={{
      margin: '100px',
    }}
  >
    <Modal title={testTitle}>
      <p>This is a basic modal!</p>
    </Modal>
  </div>
);

export const withActions = () => (
  <div
    style={{
      margin: '100px',
    }}
  >
    <Modal actions={testActions}>
      <p>This is a basic modal!</p>
    </Modal>
  </div>
);

export const withTitleAndActions = () => (
  <div
    style={{
      margin: '100px',
    }}
  >
    <Modal title={testTitle} actions={testActions}>
      <p>This is a basic modal!</p>
    </Modal>
  </div>
);

export const closeCaseModal = () => (
  <div
    style={{
      margin: '100px',
    }}
  >
    <Modal
      title={<h2>Close case?</h2>}
      actions={
        <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
          <Button type="button" secondary>
            Cancel
          </Button>
          <Button type="button">Close investigation</Button>
        </div>
      }
    >
      <p>
        This will close your investigative work for <strong>Case #590381450</strong> and forward the case to
        adjudication.
      </p>
    </Modal>
  </div>
);

export const withOverlay = () => (
  <div>
    <Overlay />
    <ModalContainer>
      <Modal title={<h2>Close case?</h2>} actions={closeCaseActions}>
        <p>
          This will close your investigative work for <strong>Case #590381450</strong> and forward the case to
          adjudication.
        </p>
      </Modal>
    </ModalContainer>
  </div>
);

const TestModal = (closeModal) => (
  <Modal
    title={<h2>Test Modal</h2>}
    actions={
      <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
        <Button type="button" secondary onClick={closeModal}>
          Cancel
        </Button>
        <Button type="button" onClick={closeModal}>
          Close
        </Button>
      </div>
    }
  >
    <p>This is a test modal!</p>
  </Modal>
);

export const FullExample = () => {
  const { isOpen, openModal, closeModal } = useModal();

  const ConnectedTestModal = connectModal(TestModal(closeModal));

  return (
    <div>
      <Button type="button" onClick={openModal}>
        Click me!
      </Button>
      <ConnectedTestModal isOpen={isOpen} onClose={closeModal} />
    </div>
  );
};
