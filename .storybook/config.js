/* eslint-disable import/no-extraneous-dependencies */
import { configure, addDecorator } from '@storybook/react';
import { withInfo } from '@storybook/addon-info';
import 'happo-plugin-storybook/register';

import './storybook.scss';
import '../src/index.scss';
import '../src/ghc_index.scss';

configure(require.context('../src', true, /\.stories\.jsx?$/), module);

addDecorator(withInfo);
