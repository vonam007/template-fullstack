import React from 'react'
import { useTranslation } from 'react-i18next'

const LanguageSelector: React.FC = () => {
  const { i18n, t } = useTranslation()

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng)
    localStorage.setItem('language', lng)
  }

  return (
    <div className="language-selector">
      <select 
        value={i18n.language} 
        onChange={(e) => changeLanguage(e.target.value)}
        className="form-control"
      >
        <option value="en">{t('language.english')}</option>
        <option value="vi">{t('language.vietnamese')}</option>
      </select>
    </div>
  )
}

export default LanguageSelector
