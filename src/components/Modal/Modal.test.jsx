import React from 'react';
import { shallow, mount } from 'enzyme';

import Modal, { ModalTitle, ModalClose, ModalActions, connectModal } from './Modal';

describe('Modal component', () => {
  const wrapper = mount(<Modal>Test Modal</Modal>);

  it('renders without crashing', () => {
    expect(wrapper.exists()).toBe(true);
    expect(wrapper.find('[data-testid="modal"]').exists()).toBe(true);
  });

  it('renders its children', () => {
    expect(wrapper.text()).toEqual('Test Modal');
  });
});

describe('ModalTitle component', () => {
  const wrapper = shallow(<ModalTitle>Test Modal Title</ModalTitle>);

  it('renders its children', () => {
    expect(wrapper.text()).toEqual('Test Modal Title');
  });
});

describe('ModalActions component', () => {
  const wrapper = shallow(<ModalActions>Close</ModalActions>);

  it('renders its children', () => {
    expect(wrapper.text()).toEqual('Close');
  });
});

describe('ModalClose component', () => {
  const mockClose = jest.fn();
  const wrapper = mount(<ModalClose handleClick={mockClose} />);

  it('renders a close button', () => {
    expect(wrapper.find('button[data-testid="modalCloseButton"]').exists()).toBe(true);
  });

  it('calls the close handler onClick', () => {
    wrapper.simulate('click');
    expect(mockClose).toHaveBeenCalled();
  });
});

describe('connectModal HOC', () => {
  const MockModal = () => <Modal>Test Connected Modal</Modal>;
  const ConnectedModal = connectModal(MockModal);

  it('renders without crashing', () => {
    const wrapper = mount(<ConnectedModal />);
    expect(wrapper.find('ConnectedMockModal').exists()).toBe(true);
  });
});
