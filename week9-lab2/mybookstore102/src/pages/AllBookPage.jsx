import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { BookOpenIcon, LogoutIcon } from '@heroicons/react/outline';

const AllBookPage = () => {
    const [Books, setBooks] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchBooks = async () => {
            try {
                setLoading(true);
                const response = await fetch('/api/v1/books/');
                if (!response.ok) {
                    throw new Error('Failed to fetch books');
                }
                const data = await response.json();
                setBooks(data);
                setError(null);
            } catch (err) {
                setError(err.message);
                console.error('Error fetching books:', err);
            } finally {
                setLoading(false);
            }
        };

        fetchBooks();
    }, []);

    const handleEdit = (book) => {
        console.log('Edit book:', book);
        alert(`Edit: ${book.title}`);
        // navigate(`/books/${book.id}/edit`);
    };

    const handleDelete = async (id) => {
        const ok = window.confirm('ต้องการลบหนังสือเล่มนี้หรือไม่?');
        if (!ok) return;

        try {
            const res = await fetch(`/api/v1/books/${id}`, { method: 'DELETE' });
            if (!res.ok) {
                throw new Error('ลบไม่สำเร็จ');
            }
            setBooks((prev) => prev.filter((b) => b.id !== id));
        } catch (err) {
            console.error(err);
            alert(err.message || 'เกิดข้อผิดพลาดระหว่างลบข้อมูล');
        }
    };

    const handleLogout = () => {
        localStorage.removeItem('isAdminAuthenticated');
        navigate('/login');
    };

    // Loading & Error — ให้ดูเป็นกลาง ๆ ตรงกลางหน้า
    if (loading) {
        return (
            <div className="min-h-screen bg-gray-50">
                <header className="bg-gradient-to-r from-green-600 to-green-700 text-white shadow-lg">
                    <div className="container mx-auto px-4 py-6">
                        <div className="flex justify-between items-center">
                            <div className="flex items-center space-x-3">
                                <BookOpenIcon className="h-8 w-8" />
                                <h1 className="text-2xl font-bold">BookStore - BackOffice</h1>
                            </div>
                            <button
                                onClick={handleLogout}
                                className="flex items-center space-x-2 px-4 py-2 bg-white/20 hover:bg-white/30 rounded-lg transition-colors"
                            >
                                <LogoutIcon className="h-5 w-5" />
                                <span>ออกจากระบบ</span>
                            </button>
                        </div>
                    </div>
                </header>
                <main className="container mx-auto px-4 py-10">
                    <div className="max-w-5xl mx-auto text-center text-gray-700">Loading...</div>
                </main>
            </div>
        );
    }

    if (error) {
        return (
            <div className="min-h-screen bg-gray-50">
                <header className="bg-gradient-to-r from-green-600 to-green-700 text-white shadow-lg">
                    <div className="container mx-auto px-4 py-6">
                        <div className="flex justify-between items-center">
                            <div className="flex items-center space-x-3">
                                <BookOpenIcon className="h-8 w-8" />
                                <h1 className="text-2xl font-bold">BookStore - BackOffice</h1>
                            </div>
                            <button
                                onClick={handleLogout}
                                className="flex items-center space-x-2 px-4 py-2 bg-white/20 hover:bg-white/30 rounded-lg transition-colors"
                            >
                                <LogoutIcon className="h-5 w-5" />
                                <span>ออกจากระบบ</span>
                            </button>
                        </div>
                    </div>
                </header>
                <main className="container mx-auto px-4 py-10">
                    <div className="max-w-5xl mx-auto text-center text-red-600">Error: {error}</div>
                </main>
            </div>
        );
    }

    // กรณีข้อมูลว่าง
    if (!Books || Books.length === 0) {
        return (
            <div className="min-h-screen bg-gray-50">
                <header className="bg-gradient-to-r from-green-600 to-green-700 text-white shadow-lg">
                    <div className="container mx-auto px-4 py-6">
                        <div className="flex justify-between items-center">
                            <div className="flex items-center space-x-3">
                                <BookOpenIcon className="h-8 w-8" />
                                <h1 className="text-2xl font-bold">BookStore - BackOffice</h1>
                            </div>
                            <button
                                onClick={handleLogout}
                                className="flex items-center space-x-2 px-4 py-2 bg-white/20 hover:bg-white/30 rounded-lg transition-colors"
                            >
                                <LogoutIcon className="h-5 w-5" />
                                <span>ออกจากระบบ</span>
                            </button>
                        </div>
                    </div>
                </header>
                <main className="container mx-auto px-4 py-10">
                    <div className="max-w-5xl mx-auto text-center text-gray-600">ยังไม่มีรายการหนังสือ</div>
                </main>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-gray-50">
            {/* Header เต็มกว้างเหมือนเดิม */}
            <header className="bg-gradient-to-r from-green-600 to-green-700 text-white shadow-lg">
                <div className="container mx-auto px-4 py-6">
                    <div className="flex justify-between items-center">
                        <div className="flex items-center space-x-3">
                            <BookOpenIcon className="h-8 w-8" />
                            <h1 className="text-2xl font-bold">BookStore - BackOffice</h1>
                        </div>
                        <button
                            onClick={handleLogout}
                            className="flex items-center space-x-2 px-4 py-2 bg-white/20 hover:bg-white/30 rounded-lg transition-colors"
                        >
                            <LogoutIcon className="h-5 w-5" />
                            <span>ออกจากระบบ</span>
                        </button>
                    </div>
                </div>
            </header>
            {/* ส่วนเนื้อหา: จำกัดความกว้างและจัดกลาง */}
            <main className="container mx-auto px-4 py-10">
                <div className="w-full max-w-5xl mx-auto">
                <h1>จัดการหนังสือทิ้งหมด</h1>
                    <div className="overflow-x-auto rounded-lg shadow border border-gray-200 bg-white">
                        <table className="w-full">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border-b">ID</th>
                                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border-b">Title</th>
                                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border-b">Author</th>
                                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border-b">ISBN</th>
                                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border-b">Year</th>
                                    <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700 border-b">Price</th>
                                    <th className="px-4 py-3 text-right text-sm font-semibold text-gray-700 border-b">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {Books.map((book) => (
                                    <tr key={book.id} className="hover:bg-gray-50">
                                        <td className="px-4 py-3 border-b text-sm text-gray-800">{book.id}</td>
                                        <td className="px-4 py-3 border-b text-sm text-gray-800">{book.title}</td>
                                        <td className="px-4 py-3 border-b text-sm text-gray-800">{book.author}</td>
                                        <td className="px-4 py-3 border-b text-sm text-gray-800">{book.isbn}</td>
                                        <td className="px-4 py-3 border-b text-sm text-gray-800">{book.year}</td>
                                        <td className="px-4 py-3 border-b text-sm text-gray-800">
                                            {typeof book.price === 'number' ? book.price.toFixed(2) : book.price}
                                        </td>
                                        <td className="px-4 py-3 border-b">
                                            <div className="flex items-center justify-end gap-2">
                                                <button
                                                    onClick={() => handleEdit(book)}
                                                    className="px-3 py-1.5 rounded-lg text-sm border border-blue-500 text-blue-600 hover:bg-blue-50"
                                                >
                                                    Edit
                                                </button>
                                                <button
                                                    onClick={() => handleDelete(book.id)}
                                                    className="px-3 py-1.5 rounded-lg text-sm border border-red-500 text-red-600 hover:bg-red-50"
                                                >
                                                    Delete
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                    {/* เพิ่มระยะด้านล่างเล็กน้อย */}
                    <div className="h-6" />
                </div>
            </main>
        </div>
    );
};

export default AllBookPage;
