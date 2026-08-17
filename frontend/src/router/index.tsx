import { createBrowserRouter } from 'react-router-dom';
import Shell from '../components/Shell';
import Login from '../pages/Login';
import Home from '../pages/Home';
import ArticleDetail from '../pages/ArticleDetail';
import CategoryArticles from '../pages/CategoryArticles';
import TagArticles from '../pages/TagArticles';
import Archives from '../pages/Archives';
import Search from '../pages/Search';
import Admin from '../pages/Admin';
import { RequireAdmin } from './guards';

export const router = createBrowserRouter([
  { path: '/login', element: <Login /> },
  {
    path: '/',
    element: <Shell />,
    children: [
      { index: true, element: <Home /> },
      { path: 'articles/:slug', element: <ArticleDetail /> },
      { path: 'categories/:slug', element: <CategoryArticles /> },
      { path: 'tags/:slug', element: <TagArticles /> },
      { path: 'archives', element: <Archives /> },
      { path: 'search', element: <Search /> },
      { path: 'admin', element: <RequireAdmin><Admin /></RequireAdmin> },
    ],
  },
]);
